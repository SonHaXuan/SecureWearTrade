// Simple deployment script without confirmations
const hre = require("hardhat");

async function deployToNetwork(networkName) {
  console.log(`\n=== Deploying to ${networkName} ===`);

  try {
    console.log(`Deploying KYC contract to ${networkName}...`);
    const KYC = await hre.ethers.getContractFactory("KYC");
    const kyc = await KYC.deploy();

    await kyc.deployed();

    console.log(`✅ KYC contract deployed to: ${kyc.address}`);
    console.log(`📝 Transaction hash: ${kyc.deployTransaction.hash}`);

    // Return deployment info
    return {
      network: networkName,
      address: kyc.address,
      txHash: kyc.deployTransaction.hash,
      explorer: getExplorerUrl(networkName, kyc.address)
    };

  } catch (error) {
    console.error(`❌ Deployment failed on ${networkName}:`, error.message);
    return null;
  }
}

function getExplorerUrl(networkName, address) {
  const explorers = {
    sepolia: `https://sepolia.etherscan.io/address/${address}`,
    bscTestnet: `https://testnet.bscscan.com/address/${address}`,
    phantomTestnet: `https://testnet.ftmscan.com/address/${address}`,
    celoTestnet: `https://alfajores-blockscout.celo-testnet.org/address/${address}`,
    polygonTestnet: `https://mumbai.polygonscan.com/address/${address}`
  };
  return explorers[networkName] || '';
}

async function main() {
  const networkName = hre.network.name;
  console.log(`🚀 Deploying KYC Contract to ${networkName}`);

  const deployment = await deployToNetwork(networkName);

  if (deployment) {
    console.log(`\n🎉 Deployment successful!`);
    console.log(`📍 Network: ${deployment.network}`);
    console.log(`🏠 Contract: ${deployment.address}`);
    console.log(`🔍 Explorer: ${deployment.explorer}`);
    console.log(`📝 Transaction: ${deployment.txHash}`);
  } else {
    console.log(`❌ Deployment failed. Please check your wallet balance and network configuration.`);
    console.log(`💰 Fund your wallet with testnet tokens from the respective faucets.`);
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});