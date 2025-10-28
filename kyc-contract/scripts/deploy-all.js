// Multi-network deployment script with balance checking
const hre = require("hardhat");

async function checkBalance(networkName, provider) {
  try {
    const [signer] = await ethers.getSigners();
    const balance = await signer.getBalance();
    const formattedBalance = ethers.utils.formatEther(balance);
    console.log(`${networkName} balance: ${formattedBalance} ETH`);
    return balance;
  } catch (error) {
    console.error(`Error checking ${networkName} balance:`, error.message);
    return ethers.BigNumber.from(0);
  }
}

async function deployToNetwork(networkName) {
  console.log(`\n=== Deploying to ${networkName} ===`);

  try {
    // Get network provider
    const network = hre.config.networks[networkName];
    if (!network) {
      throw new Error(`Network ${networkName} not configured`);
    }

    // Check balance first
    const provider = new ethers.providers.JsonRpcProvider(network.url);
    const balance = await checkBalance(networkName, provider);

    if (balance.eq(0)) {
      console.log(`⚠️  Insufficient funds on ${networkName}. Please fund the wallet first.`);
      console.log(`💰 Required: ~0.01 ${networkName.includes('bsc') || networkName.includes('phantom') ? 'BNB/FTM' : networkName.includes('celo') ? 'CELO' : 'MATIC'}`);
      console.log(`🔗 Wallet: ${provider.getSigner().getAddress()}`);
      return null;
    }

    // Deploy contract
    console.log(`Deploying KYC contract to ${networkName}...`);
    const KYC = await hre.ethers.getContractFactory("KYC");
    const kyc = await KYC.deploy();

    await kyc.deployed();

    console.log(`✅ KYC contract deployed to: ${kyc.address}`);
    console.log(`📝 Transaction hash: ${kyc.deployTransaction.hash}`);

    // Wait for confirmations
    console.log("⏳ Waiting for block confirmations...");
    await kyc.deployTransaction.wait(2);
    console.log("✅ Contract deployment confirmed!");

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
  console.log("🚀 Multi-Network KYC Contract Deployment");
  console.log("========================================");

  const networks = ['bscTestnet', 'phantomTestnet', 'celoTestnet', 'polygonTestnet'];
  const deployments = [];

  for (const network of networks) {
    const deployment = await deployToNetwork(network);
    if (deployment) {
      deployments.push(deployment);
    }

    // Add delay between deployments
    await new Promise(resolve => setTimeout(resolve, 2000));
  }

  // Summary
  console.log("\n📊 Deployment Summary");
  console.log("====================");

  if (deployments.length === 0) {
    console.log("❌ No successful deployments. Please fund your wallet on the desired testnets.");
  } else {
    console.log(`✅ Successfully deployed to ${deployments.length} network(s):`);
    deployments.forEach(d => {
      console.log(`  📍 ${d.network}: ${d.address}`);
      console.log(`     🔍 Explorer: ${d.explorer}`);
      console.log(`     📝 Tx: ${d.txHash}`);
      console.log('');
    });
  }

  console.log("💡 Funding Requirements:");
  console.log("   - BNB Testnet: Get BNB from https://testnet.binance.org/faucet-smart");
  console.log("   - Phantom Testnet: Get FTM from https://faucet.fantom.network/");
  console.log("   - Celo Testnet: Get CELO from https://faucet.celo.org/");
  console.log("   - Polygon Testnet: Get MATIC from https://faucet.polygon.technology/");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});