import grpc from "k6/net/grpc";
import { check, sleep } from "k6";

// 1. Initialize the gRPC client
const client = new grpc.Client();

// 2. Load the proto definitions at the global scope
// Note: Replace './proto' with the actual path to your petstore.proto file
client.load(["../protobuffer/swaggerpetstore"], "petstore.proto");

export const options = {
  stages: [
    { duration: "30s", target: 20 }, // Ramp-up to 5 virtual users
    { duration: "5m", target: 50 }, // Stay at 5 VUs
    { duration: "10m", target: 100 }, // Ramp-down
  ],
};

export default () => {
  // 3. Connect to the gRPC server (use plaintext: true if not using TLS)
  client.connect("localhost:8090", { plaintext: true });

  // --- Test Case: AddPet ---
  const addPetPayload = {
    body: {
      id: 101,
      name: "Rex",
      status: 0, // STATUS_AVAILABLE
      category: { id: 1, name: "Dogs" },
      photoUrls: ["http://example.com/rex.jpg"],
      tags: [{ id: 1, name: "friendly" }],
    },
  };

  const addResponse = client.invoke(
    "petstore.SwaggerPetstoreService/AddPet",
    addPetPayload,
  );

  check(addResponse, {
    "AddPet status is OK": (r) => r && r.status === grpc.StatusOK,
  });

  // 4. Close the connection and rest
  client.close();
  sleep(1);
};
