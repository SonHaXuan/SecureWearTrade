// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// You can also run a script with `npx hardhat run <script>`. If you do that, Hardhat
// will compile your contracts, add the Hardhat Runtime Environment's members to the
// global scope, and execute the script.
const hre = require("hardhat");

async function main() {
  console.log("Deploying KYC contract to Sepolia testnet...");

  const KYC = await hre.ethers.getContractFactory("KYC");
  const kyc = await KYC.deploy();

  await kyc.deployed();

  console.log(`KYC contract deployed to: ${kyc.address}`);
  console.log(`Transaction hash: ${kyc.deployTransaction.hash}`);

  // Wait for a few block confirmations
  console.log("Waiting for block confirmations...");
  await kyc.deployTransaction.wait(2);

  console.log("Contract deployment confirmed!");
  console.log(`You can verify the contract on Etherscan: https://sepolia.etherscan.io/address/${kyc.address}`);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
