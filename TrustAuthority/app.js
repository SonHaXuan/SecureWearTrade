const express = require("express");
const chalk = require("chalk");
const app = express();
const bodyParser = require("body-parser");
const axios = require("axios");

// parse application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: false }));

// parse application/json
app.use(bodyParser.json());

const JediInstance = axios.default.create({
  baseURL: "http://localhost:8080",
});

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

app.post("/api/request-access", async (req, res) => {
  const user = req.body;

  console.log(
    chalk.blue("INFO: ") +
      chalk.bold(
        `User "${user.id}" register for access to Owner ${"Alice"}'s data\n`
      ) +
      chalk.cyan("  - URI: ") +
      `${user.requestURI}\n` +
      chalk.cyan("  - Access Expiration: ") +
      `${user.expiredDate}\n`
  );

  const jediPKResA = await JediInstance.get(`/jedi-private-key/`);
  const jediKeys = jediPKResA.data.data;

  console.log(
    chalk.blue("INFO: ") +
      chalk.bold("Key generated for access\n") +
      chalk.cyan("  - URI: ") +
      `${user.requestURI}`
  );

  console.log(
    chalk.blue("INFO: ") +
      chalk.bold("Key deletion completed\n") +
      chalk.cyan("  - New Key ID: ") +
      `${jediKeys.slice(0, 400)}...`
  );
  res.status(200).json({ key: jediKeys });
});

app.post("/api/user-data", async (req, res) => {
  const user = req.body;

  console.log(
    chalk.blue("INFO: ") +
      chalk.bold(`User "${user.id}" requested access\n`) +
      chalk.cyan("   - URI: ") +
      user.requestURI
  );

  await sleep(3000);
  // Log verification step
  console.log(
    chalk.blue("INFO: ") +
      "Verifying user access...\n" +
      chalk.cyan("   - Expiration Date Check: ") +
      chalk.gray("(in progress)")
  );

  await sleep(2000);
  console.log(
    chalk.blue("INFO: ") +
      "Access verified\n" +
      chalk.cyan("   - User'URI: ") +
      `${user.requestURI}\n` +
      chalk.cyan("   - User'Expiration date: ") +
      `${"2024-08-25T23:59:59Z"}\n` +
      chalk.cyan("   - Access Granted: ") +
      "true\n"
  );

  await sleep(5000);
  // Log data loading with shortened URLs
  console.log(
    chalk.blue("INFO: ") +
      "Loading owner’s data..." +
      chalk.gray("\n   - Files:\n") +
      chalk.cyan("      1. ") +
      "https://maroon-wandering-fly-487.mypinata.cloud/ipfs/QmaoLUdG4A7KQbbMwUN7uwpCBESaeKeHMacPhPf6Kgjk8c\n" +
      chalk.cyan("      2. ") +
      "https://maroon-wandering-fly-487.mypinata.cloud/ipfs/QmYCZziBBWp7iBBNhbiUNyqMpaHxte7qMxACvDwv7DKfvs\n" +
      chalk.cyan("      3. ") +
      "https://maroon-wandering-fly-487.mypinata.cloud/ipfs/QmTmKn8obS3DERsN74qvpiiReeF2cBE3g5CVEoRxkmToKF\n" +
      chalk.cyan("      4. ") +
      "https://maroon-wandering-fly-487.mypinata.cloud/ipfs/QmVyW1YUxyzZsgj3Zu1sLiaY6g5XaJu66dd2SpxZpak1Pi\n" +
      chalk.cyan("      5. ") +
      "https://maroon-wandering-fly-487.mypinata.cloud/ipfs/QmXwz97LD1gZuQsZZNCvj8btjbNz2wTemvdeJK6YheDyhx"
  );
  console.log(".....");

  console.log(
    chalk.blue("INFO: ") +
      "Data transfered\n" +
      chalk.cyan("   - Data Integrity Check: ") +
      "true"
  );
  console.log(chalk.blue("INFO: ") + "End of data transfer\n");

  res.status(200).json({});
});

app.post("/api/user-data2", async (req, res) => {
  const user = req.body;

  // Log the access request
  console.log(
    chalk.blue("INFO: ") +
      chalk.bold(`User "${user.id}" requested access\n`) +
      chalk.cyan("   - URI: ") +
      "/company/comsumer1"
  );

  await sleep(3000);
  // Log verification step
  console.log(
    chalk.blue("INFO: ") +
      "Verifying user access...\n" +
      chalk.cyan("   - Expiration Date Check: ") +
      chalk.gray("(in progress)")
  );

  await sleep(2000);
  console.log(
    chalk.blue("INFO: ") +
      "Access verified\n" +
      chalk.cyan("   - User'URI: ") +
      `${user.requestURI}\n` +
      chalk.cyan("   - User'Expiration date: ") +
      `${user.expiredDate}\n` +
      chalk.cyan("   - Access Granted: ") +
      "false\n" +
      chalk.cyan("   - Reason: ") +
      "User does not have access to the requested URI. Access denied."
  );
  console.log(chalk.blue("INFO: ") + "End of data transfer\n");

  res.status(200).json({});
});

app.post("/api/user-data3", async (req, res) => {
  const user = req.body;

  // Log the access request
  console.log(
    chalk.blue("INFO: ") +
      chalk.bold(`User "${user.id}" requested access\n`) +
      chalk.cyan("   - URI: ") +
      user.requestURI
  );

  await sleep(3000);
  // Log verification step
  console.log(
    chalk.blue("INFO: ") +
      "Verifying user access...\n" +
      chalk.cyan("   - Expiration Date Check: ") +
      chalk.gray("(in progress)")
  );

  await sleep(2000);
  console.log(
    chalk.blue("INFO: ") +
      "Access verified\n" +
      chalk.cyan("   - User'URI: ") +
      `${user.requestURI}\n` +
      chalk.cyan("   - User'Expiration date: ") +
      `${"2024-08-25T23:59:59Z"}\n` +
      chalk.cyan("   - Access Granted: ") +
      "false\n" +
      chalk.cyan("   - Reason: ") +
      "User's access has expired. Access denied."
  );
  console.log(chalk.blue("INFO: ") + "End of data transfer\n");

  res.status(200).json({});
});

const PORT = 3000;
console.log(`Server running on port ${PORT}`);
app.listen(PORT);

module.exports = app;
