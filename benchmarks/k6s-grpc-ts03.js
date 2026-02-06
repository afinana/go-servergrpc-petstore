import grpc from "k6/net/grpc";
import { check, sleep } from "k6";

// 1. Initialize the gRPC client
const client = new grpc.Client();

// 2. Load the proto definitions at the global scope
// Note: Replace './proto' with the actual path to your petstore.proto file
client.load([".,/protobuffer/swaggerpetstore"], "petstore.proto");

export const options = {
  stages: [
    { duration: "30s", target: 5 }, // Ramp-up to 5 virtual users
    { duration: "1m", target: 5 }, // Stay at 5 VUs
    { duration: "30s", target: 0 }, // Ramp-down
  ],
};

export default () => {
  // 3. Connect to the gRPC server (use plaintext: true if not using TLS)
  client.connect("localhost:8090", { plaintext: true });

  // --- Test Case: LoginUser ---
  const loginResponse = client.invoke(
    "petstore.SwaggerPetstoreService/LoginUser",
    {
      username: "user1",
      password: "password123",
    },
  );

  check(loginResponse, {
    "LoginUser status is OK": (r) => r && r.status === grpc.StatusOK,
  });

  // 4. Close the connection and rest
  client.close();
  sleep(1);
};
