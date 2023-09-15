import http from 'k6/http';
import { check } from 'k6';

// export let options = {
//   vus: 10, // Number of virtual users (simulated users)
//   duration: '10s', // Test duration in seconds
// };

export default function () {
  // Define the API endpoint
  const url = 'http://localhost:4000/tokens';

  // Send a POST request to create a token
  const payload = JSON.stringify({ token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c' });
  const params = {
    headers: {
      'Authorization': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c',
    },
  };
  const res = http.post(url, payload, params);

  // Check if the response status code is 200
  check(res, {
    'Status is 200': (r) => r.status === 200,
  });

  // Sleep for a short period (e.g., 1 second) between requests
//   sleep(1);
}
