// k6 run script.js

import http from 'k6/http';
import { check, sleep } from 'k6';
export const options = {
  discardResponseBodies: true,
  scenarios: {
    contacts: {
      executor: 'constant-arrival-rate',
      rate: 150,
      timeUnit: '1s',
      duration: '20s',
      preAllocatedVUs: 50,
      maxVUs: 1200,
    },
  },
};
// test HTTP
export default function () {
  const res = http.get('http://localhost:5000/method1');
  check(res, { 'status was 200': (r) => r.status == 200 });
  sleep(1);
}