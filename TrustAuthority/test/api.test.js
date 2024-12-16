const request = require("supertest");
const app = require("../app");

var jediKeys;
describe("API Endpoints", () => {
  test(
    "Should request access to Onwer successfully",
    async () => {
      const response = await request(app)
        .post("/api/request-access")
        .send({
          id: "company1",
          requestURI: "company/comsumer1/comsumer2",
          expiredDate: "2024-11-21T12:00:00Z",
        })
        .set("Content-Type", "application/json");

      expect(response.statusCode).toBe(200);
      expect(typeof response.body.key).toBe("string");

      jediKeys = response.body.key;
    },
    30 * 10_000
  );

  test(
    ">>>> Case 1: User Access Authorized...",
    async () => {
      const response = await request(app)
        .post("/api/user-data")
        .send({
          id: "company1",
          requestURI: "company/comsumer1/comsumer2",
          expiredDate: "2024-11-21T12:00:00Z",
        })
        .set("Content-Type", "application/json");

      expect(response.statusCode).toBe(200);
    },
    30 * 10_000
  );

  test(
    ">>>> Case 2: User Access Denied...",
    async () => {
      const response = await request(app)
        .post("/api/user-data2")
        .send({
          id: "company1",
          requestURI: "company/comsumer1/comsumer2",
          expiredDate: "2024-11-21T12:00:00Z",
        })
        .set("Content-Type", "application/json");

      expect(response.statusCode).toBe(200);
    },
    30 * 10_000
  );

  test(
    ">>>> Case 3: User Access Expired...",
    async () => {
      const response = await request(app)
        .post("/api/user-data3")
        .send({
          id: "company1",
          requestURI: "company/comsumer1/comsumer2",
          expiredDate: "2024-11-21T12:00:00Z",
        })
        .set("Content-Type", "application/json");

      expect(response.statusCode).toBe(200);
    },
    30 * 10_000
  );
});
