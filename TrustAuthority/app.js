const express = require("express");
const chalk = require("chalk");
const app = express();
const bodyParser = require("body-parser");
const axios = require("axios");
const path = require("path");
const fs = require("fs");
const moment = require("moment");

// parse application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: false }));

// parse application/json
app.use(bodyParser.json());

const JediInstance = axios.default.create({
  baseURL: "http://localhost:8080",
});

// Path to the JSON file
const userFilePath = path.join(__dirname, "users.json");
const userData = JSON.parse(
  fs.readFileSync(path.join(__dirname, "user-data.json"), "utf8")
);

// Function to initialize the JSON file
function initializeFile() {
  if (!fs.existsSync(userFilePath)) {
    // If file does not exist, create it with an empty array
    fs.writeFileSync(userFilePath, JSON.stringify([]), "utf8");
    console.log("File created successfully");
  } else {
    console.log("File already exists");
  }
}
initializeFile();

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

app.post("/api/request-access", async (req, res) => {
  const user = req.body;

  console.log(
    chalk.blue("INFO: ") +
      chalk.bold(
        `User "${user.id}" register for access to Owner ${userData.name}'s data\n`
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

  const data = JSON.parse(fs.readFileSync(userFilePath, "utf8"));
  data.push(user); // Add new user to the array
  fs.writeFileSync(userFilePath, JSON.stringify(data, null, 2), "utf8"); // Save back to file

  res.status(200).json({ key: jediKeys });
});

app.post("/api/user-data", async (req, res) => {
  const user = req.body;
  const data = JSON.parse(fs.readFileSync(userFilePath, "utf8"));

  const existedUser = data.find((u) => u.id === user.id);

  console.log(
    chalk.blue("INFO: ") +
      chalk.bold(`User "${user.id}" requested access\n`) +
      chalk.cyan("   - URI: ") +
      user.requestURI
  );

  // Log verification step
  console.log(
    chalk.blue("INFO: ") +
      "Verifying user access...\n" +
      chalk.cyan("   - Expiration Date Check: ") +
      chalk.gray("(in progress)")
  );

  if (moment().isAfter(moment(existedUser.expiredDate))) {
    console.log(
      chalk.blue("INFO: ") +
        "Access verified\n" +
        chalk.cyan("   - User'URI: ") +
        `${existedUser.requestURI}\n` +
        chalk.cyan("   - User'Expiration date: ") +
        `${existedUser.expiredDate}\n` +
        chalk.cyan("   - Access Granted: ") +
        "false\n" +
        chalk.cyan("   - Reason: ") +
        "User's access has expired. Access denied."
    );
    console.log(chalk.blue("INFO: ") + "End of data transfer\n");

    return res.status(200).json({});
  } else if (existedUser.requestURI !== user.requestURI) {
    console.log(
      chalk.blue("INFO: ") +
        "Access verified\n" +
        chalk.cyan("   - User'URI: ") +
        `${existedUser.requestURI}\n` +
        chalk.cyan("   - User'Expiration date: ") +
        `${existedUser.expiredDate}\n` +
        chalk.cyan("   - Access Granted: ") +
        "false\n" +
        chalk.cyan("   - Reason: ") +
        "User does not have access to the requested URI. Access denied."
    );
    console.log(chalk.blue("INFO: ") + "End of data transfer\n");

    return res.status(200).json({});
  }

  console.log(
    chalk.blue("INFO: ") +
      "Access verified\n" +
      chalk.cyan("   - User'URI: ") +
      `${user.requestURI}\n` +
      chalk.cyan("   - User'Expiration date: ") +
      `${existedUser.expiredDate}\n` +
      chalk.cyan("   - Access Granted: ") +
      "true\n"
  );

  await sleep(5000);
  // Log data loading with shortened URLs
  console.log(
    chalk.blue("INFO: ") +
      "Loading owner’s data..." +
      chalk.gray("\n   - Files:\n") +
      `${userData.dataPoints
        .map((dataPoint, index) => {
          return chalk.cyan(`      ${index + 1}. `) + `${dataPoint.url}`;
        })
        .join("\n")}`
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
