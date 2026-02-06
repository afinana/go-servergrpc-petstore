To test your gRPC API using Grafana k6, you will need to use the `k6/net/grpc` module. Since k6 uses JavaScript but the gRPC service communicates via Protocol Buffers, the test script must load your `.proto` file to understand the service definitions and message structures.

### 1. Prerequisites

* **Install k6:** Ensure you have the latest version of k6 installed (v0.29.0 or higher is required for native gRPC support).
* **Proto File:** Place your `petstore.proto` in a directory accessible to your script (e.g., a folder named `/../protobuffer/swaggerpetstore`).

---

### 2. k6 Test Script Example


- `k6s-grpc-ts01.js` is an example script that demonstrates how to load your proto file and make a gRPC call to the `addPetPayload` method of the `PetService` API.
- `k6s-grpc-ts02.js` is an example script that demonstrates how to load your proto file and make a gRPC call to the `GetPetById` method of the `PetService` API with plaintext connection.
- `k6s-grpc-ts02.js` is an example script that demonstrates how to load your proto file and make a gRPC call to the `LoginUser` method of the API with plaintext connection.


### 3. Key Components Explained

* **`client.load(importPaths, protoFile)`**: This loads your definitions. The first argument is an array of paths to search for imports (like `google/protobuf/empty.proto` used in your file), and the second is the filename.
* **`client.invoke(url, data)`**: The URL follows the format `package.ServiceName/MethodName`. In your case, the package is `petstore`.
* **Status Checks**: Use `grpc.StatusOK` to verify that the call was successful. You can also inspect `response.message` for the data returned by the server.
* **Plaintext vs. TLS**: By default, k6 expects an encrypted connection. If your local dev server does not use SSL/TLS, you must set `{ plaintext: true }` in the `connect` options.

### 4. Running the Test

Execute the test from your terminal:

```bash
k6 run your_script_name.js
```

To configure load scenarios for your gRPC Petstore service, you must define the `options` object in your k6 script to control the number of Virtual Users (VUs) and the duration of the test.

Each scenario below targets a different performance goal, such as finding the breaking point of the service or ensuring it recovers after a sudden burst of traffic.

---

## 4. Stress Testing

A **Stress Test** determines the maximum capacity of your service by gradually increasing the load until the system begins to fail or performance degrades significantly.

```javascript
export const options = {
  stages: [
    { duration: '2m', target: 50 },  // Ramp-up: gradual increase to 50 users
    { duration: '5m', target: 50 },  // Plateau: stay at 50 users to find bottlenecks
    { duration: '2m', target: 100 }, // Ramp-up: push to 100 users
    { duration: '5m', target: 100 }, // Plateau: stay at 100 users
    { duration: '2m', target: 0 },   // Ramp-down: cool down
  ],
  thresholds: {
    'grpc_req_duration': ['p(95)<500'], // 95% of gRPC calls must complete under 500ms
  },
};

```

## 5. Spike Testing

A **Spike Test** simulates a sudden, massive increase in traffic (e.g., a marketing event or a flash sale) to see if the service can handle the shock and recover afterward.

```javascript
export const options = {
  stages: [
    { duration: '10s', target: 10 },  // Normal baseline
    { duration: '1m', target: 200 }, // Sudden spike: jump to 200 users instantly
    { duration: '2m', target: 200 }, // Hold the spike
    { duration: '1m', target: 10 },  // Rapid recovery: drop back to baseline
    { duration: '10s', target: 0 },
  ],
};

```

## 6. Soak (Endurance) Testing

A **Soak Test** checks for system stability over a long period. This is essential for identifying memory leaks in your gRPC server or database connection exhaustion.

```javascript
export const options = {
  stages: [
    { duration: '2m', target: 20 },  // Ramp-up
    { duration: '2h', target: 20 },  // Soak: stay at a steady load for 2 hours
    { duration: '2m', target: 0 },   // Ramp-down
  ],
};
```

---
### Implementation Tips for Petstore gRPC

* **Thresholds:** Always use `thresholds` in your k6 options to automatically fail the test if the error rate exceeds a certain percentage (e.g., `grpc_req_status: ['rate < 0.01']` ensures less than 1% failure).
* **Dynamic Data:** For these scenarios, avoid using the same `petId` or `username` repeatedly to prevent skewed results caused by database caching.
* 
**Protobuf Methods:** Since your `.proto` includes `DeleteOrder` and `DeletePet`, ensure your load tests don't delete data that other concurrent VUs are trying to "Get," unless you are specifically testing race conditions.


By implementing these load testing scenarios, you can gain insights into the performance characteristics of your gRPC Petstore service under various conditions, helping you identify bottlenecks and ensure reliability in production.
