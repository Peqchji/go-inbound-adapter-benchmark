import http from 'k6/http';
import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';

// Custom Metrics for Latency Comparison
const restLatency = new Trend('custom_rest_req_duration');
const grpcLatency = new Trend('custom_grpc_req_duration');
const gqlLatency = new Trend('custom_gql_req_duration');

const HTTP_HOST = __ENV.HTTP_HOST || 'localhost:8080';
const GQL_HOST = __ENV.GQL_HOST || 'localhost:8081';
const GRPC_HOST = __ENV.GRPC_HOST || 'localhost:50051';

const grpcClient = new grpc.Client();
grpcClient.load(['../cmd/grpcserver/proto'], 'wallet.proto');
const vus = 50

const restParams = { headers: { 'Content-Type': 'application/json' } };
const gqlHeaders = { 'Content-Type': 'application/json' };


export const options = {
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
    thresholds: {
        custom_rest_req_duration: ['p(95)<500'],
        custom_grpc_req_duration: ['p(95)<500'],
        custom_gql_req_duration: ['p(95)<500'],
    },
    scenarios: {
        rest_test: {
            executor: 'constant-vus',
            vus: vus,
            duration: '10s',
            exec: 'restBenchmark',
            startTime: '0s', // Run first
        },
        grpc_test: {
            executor: 'constant-vus',
            vus: vus,
            duration: '10s',
            exec: 'grpcBenchmark',
            startTime: '20s', // Run after HTTP (approx) - or can run parallel if servers handle it
        },
        gql_test: {
            executor: 'constant-vus',
            vus: vus,
            duration: '10s',
            exec: 'gqlBenchmark',
            startTime: '40s',
        },
    },
};

export function restBenchmark() {
    const url = `http://${HTTP_HOST}/wallets`;
    const payload = JSON.stringify({
        id: `http-${__VU}-${__ITER}`,
        firstname: "John",
        lastname: "Doe",
    });
    const start = Date.now();
    const res = http.post(url, payload, restParams);
    const duration = Date.now() - start;
    restLatency.add(duration);

    check(res, { 'http 201': (r) => r.status === 201 });

    const startGet = Date.now();
    const getRes = http.get(`${url}/http-${__VU}-${__ITER}`, { tags: { name: 'GetWalletById' } });
    const durationGet = Date.now() - startGet;
    restLatency.add(durationGet);

    check(getRes, { 'http 200': (r) => r.status === 200 });
}

let isGrpcConnected = false;

export function grpcBenchmark() {
    if (!isGrpcConnected) {
        grpcClient.connect(GRPC_HOST, { plaintext: true });
        isGrpcConnected = true;
    }

    const id = `grpc-${__VU}-${__ITER}`;
    const data = { owner_id: id, firstname: "John", lastname: "Doe" };

    const start = Date.now();
    const res = grpcClient.invoke('wallet.WalletService/CreateWallet', data);
    const duration = Date.now() - start;
    grpcLatency.add(duration);

    check(res, { 'grpc create OK': (r) => r && r.status === grpc.StatusOK });

    const getData = { id: id };
    const startGet = Date.now();
    const getRes = grpcClient.invoke('wallet.WalletService/GetWallet', getData);
    const durationGet = Date.now() - startGet;
    grpcLatency.add(durationGet);

    check(getRes, { 'grpc get OK': (r) => r && r.status === grpc.StatusOK });
}

export function gqlBenchmark() {
    const url = `http://${GQL_HOST}/query`;
    const id = `gql-${__VU}-${__ITER}`;
    const mutation = `mutation { createWallet(input: {id: "${id}", firstname: "John", lastname: "Doe"}) { id } }`;

    const start = Date.now();
    const res = http.post(url, JSON.stringify({ query: mutation }), { headers: gqlHeaders });
    const duration = Date.now() - start;
    gqlLatency.add(duration);

    check(res, { 'gql create 200': (r) => r.status === 200 });

    const query = `query { wallet(id: "${id}") { id balance } }`;
    const startGet = Date.now();
    const getRes = http.post(url, JSON.stringify({ query: query }), { headers: gqlHeaders });
    const durationGet = Date.now() - startGet;
    gqlLatency.add(durationGet);

    check(getRes, { 'gql get 200': (r) => r.status === 200 });
}
