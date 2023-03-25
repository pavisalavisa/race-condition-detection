import http from 'k6/http';
import { check } from 'k6';

const SERVICE_A_URL = 'http://localhost:8080/a/foo/bar'
const SERVICE_A_EXPECTED_RESPONSE = 'Hello from service A'
const SERVICE_A_METHOD = "POST"
const SERVICE_B_URL = 'http://localhost:8080/b/baz/14'
const SERVICE_B_EXPECTED_RESPONSE = 'Hello from service B'
const SERVICE_B_METHOD = "GET"

export default function() {
  // Randomly choose between two URLs
  const url = Math.random() > 0.5 ?  SERVICE_A_URL: SERVICE_B_URL;
  const expectedResponse = url === SERVICE_A_URL ? SERVICE_A_EXPECTED_RESPONSE : SERVICE_B_EXPECTED_RESPONSE
  const method = url === SERVICE_A_URL ? SERVICE_A_METHOD : SERVICE_B_METHOD

  // Make the GET request
  const res = http.request(method, url);

  // Check that the response was successful
  check(res, {
    'status is 200': (r) => r.status === 200,
    'OK response': (r)=> r.body == expectedResponse
  });
}
