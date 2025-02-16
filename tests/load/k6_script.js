import http from 'k6/http';
import { check } from 'k6';

export let options = {
  scenarios: {
    info_load: {
      executor: 'constant-arrival-rate',
      rate: 1000,
      timeUnit: '1s',
      duration: '1m',
      preAllocatedVUs: 200,
      maxVUs: 500,
    },
  },
};

export default function () {
  let res = http.get('http://localhost:8080/api/info');
  check(res, {
    'status is 200': (r) => r.status === 200,
  });
}
