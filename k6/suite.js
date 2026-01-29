import http from 'k6/http';
import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';

// Custom Metrics for Latency Comparison
const httpLatency = new Trend('http_duration_custom');
const grpcLatency = new Trend('grpc_duration_custom');
const gqlLatency = new Trend('gql_duration_custom');

const grpcClient = new grpc.Client();
grpcClient.load(['../internal/handler/grpc/proto'], 'wallet.proto');
const vus = 50

export const options = {
    scenarios: {
        http_test: {
            executor: 'constant-vus',
            vus: vus,
            duration: '10s',
            exec: 'httpBenchmark',
            startTime: '0s', // Run first
        },
        grpc_test: {
            executor: 'constant-vus',
            vus: vus,
            duration: '10s',
            exec: 'grpcBenchmark',
            startTime: '10s', // Run after HTTP (approx) - or can run parallel if servers handle it
        },
        gql_test: {
            executor: 'constant-vus',
            vus: vus,
            duration: '10s',
            exec: 'gqlBenchmark',
            startTime: '20s', // Run after gRPC
        },
    },
};

export function httpBenchmark() {
    const url = 'http://localhost:8080/wallets';
    const payload = JSON.stringify({
        id: `http-${__VU}-${__ITER}`,
        balance: 1000,
    });
    const params = { headers: { 'Content-Type': 'application/json' } };

    const start = Date.now();
    const res = http.post(url, payload, params);
    const duration = Date.now() - start;
    httpLatency.add(duration);

    check(res, { 'http 201': (r) => r.status === 201 });

    const getRes = http.get(`${url}/http-${__VU}-${__ITER}`);
    check(getRes, { 'http 200': (r) => r.status === 200 });
}

export function grpcBenchmark() {
    grpcClient.connect('localhost:50051', { plaintext: true });

    const data = { id: `grpc-${__VU}-${__ITER}`, balance: 1000 };

    const start = Date.now();
    const res = grpcClient.invoke('wallet.WalletService/CreateWallet', data);
    const duration = Date.now() - start;
    grpcLatency.add(duration);

    check(res, { 'grpc OK': (r) => r && r.status === grpc.StatusOK });

    grpcClient.close();
}

export function gqlBenchmark() {
    const url = 'http://localhost:8081/query';
    const headers = { 'Content-Type': 'application/json' };
    const mutation = `mutation { createWallet(input: {id: "gql-${__VU}-${__ITER}", balance: 1000}) { id } }`;

    const start = Date.now();
    const res = http.post(url, JSON.stringify({ query: mutation }), { headers: headers });
    const duration = Date.now() - start;
    gqlLatency.add(duration);

    check(res, { 'gql 200': (r) => r.status === 200 });
}
