#!/usr/bin/env node

/**
 * Daily WETH Transfer Script
 * 
 * This script transfers 2 ETH worth of WETH9 (minus gas costs) to yaketh.eth
 * on Ethereum mainnet.
 */

const { ethers } = require('ethers');
const fs = require('fs');

// WETH9 contract address on Ethereum mainnet
const WETH9_ADDRESS = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2';

// WETH9 ABI (minimal - only what we need)
const WETH9_ABI = [
  'function balanceOf(address) view returns (uint256)',
  'function withdraw(uint256) external',
  'function transfer(address to, uint256 amount) returns (bool)'
];

async function main() {
  const timestamp = new Date().toISOString();
  const logFile = `transfer-log-${timestamp.replace(/:/g, '-')}.txt`;
  
  function log(message) {
    const logMessage = `[${new Date().toISOString()}] ${message}`;
    console.log(logMessage);
    fs.appendFileSync(logFile, logMessage + '\n');
  }

  try {
    // Validate environment variables
    const privateKey = process.env.PRIVATE_KEY;
    const rpcUrl = process.env.RPC_URL;
    const recipientAddress = process.env.RECIPIENT_ADDRESS || 'yaketh.eth';
    const amountEth = process.env.AMOUNT_ETH || '2';
    const startDate = process.env.START_DATE || '2026-04-08';

    if (!privateKey) {
      throw new Error('PRIVATE_KEY environment variable is required');
    }
    if (!rpcUrl) {
      throw new Error('RPC_URL environment variable is required');
    }

    log('Starting WETH transfer process...');
    log(`Start date: ${startDate}`);
    log(`Recipient: ${recipientAddress}`);
    log(`Target amount: ${amountEth} ETH`);

    // Connect to Ethereum mainnet
    const provider = new ethers.JsonRpcProvider(rpcUrl);
    const wallet = new ethers.Wallet(privateKey, provider);
    
    log(`Wallet address: ${wallet.address}`);

    // Get network info
    const network = await provider.getNetwork();
    log(`Connected to network: ${network.name} (chainId: ${network.chainId})`);
    
    if (network.chainId !== 1n) {
      throw new Error(`Expected Ethereum mainnet (chainId 1), got chainId ${network.chainId}`);
    }

    // Resolve ENS name to address
    log(`Resolving ENS name: ${recipientAddress}`);
    const resolvedAddress = await provider.resolveName(recipientAddress);
    
    if (!resolvedAddress) {
      throw new Error(`Failed to resolve ENS name: ${recipientAddress}`);
    }
    
    log(`Resolved to address: ${resolvedAddress}`);

    // Connect to WETH9 contract
    const weth = new ethers.Contract(WETH9_ADDRESS, WETH9_ABI, wallet);
    
    // Check WETH balance
    const wethBalance = await weth.balanceOf(wallet.address);
    log(`WETH balance: ${ethers.formatEther(wethBalance)} WETH`);

    const targetAmount = ethers.parseEther(amountEth);
    
    if (wethBalance < targetAmount) {
      throw new Error(
        `Insufficient WETH balance. Have: ${ethers.formatEther(wethBalance)} WETH, Need: ${amountEth} WETH`
      );
    }

    // Estimate gas for the transfer
    log('Estimating gas for transfer...');
    const gasEstimate = await weth.transfer.estimateGas(resolvedAddress, targetAmount);
    const feeData = await provider.getFeeData();
    
    const gasPrice = feeData.gasPrice || feeData.maxFeePerGas;
    const estimatedGasCost = gasEstimate * gasPrice;
    
    log(`Estimated gas: ${gasEstimate.toString()} units`);
    log(`Gas price: ${ethers.formatUnits(gasPrice, 'gwei')} gwei`);
    log(`Estimated gas cost: ${ethers.formatEther(estimatedGasCost)} ETH`);

    // Calculate amount minus gas
    const amountMinusGas = targetAmount - estimatedGasCost;
    
    if (amountMinusGas <= 0n) {
      throw new Error('Gas cost exceeds transfer amount');
    }
    
    log(`Amount to transfer (minus gas): ${ethers.formatEther(amountMinusGas)} WETH`);

    // Check ETH balance for gas
    const ethBalance = await provider.getBalance(wallet.address);
    log(`ETH balance: ${ethers.formatEther(ethBalance)} ETH`);
    
    if (ethBalance < estimatedGasCost * 2n) { // 2x buffer for safety
      throw new Error(
        `Insufficient ETH for gas. Have: ${ethers.formatEther(ethBalance)} ETH, Need: ~${ethers.formatEther(estimatedGasCost * 2n)} ETH`
      );
    }

    // Execute the transfer
    log('Executing WETH transfer...');
    const tx = await weth.transfer(resolvedAddress, amountMinusGas);
    
    log(`Transaction submitted: ${tx.hash}`);
    log('Waiting for confirmation...');
    
    const receipt = await tx.wait();
    
    log(`Transaction confirmed in block ${receipt.blockNumber}`);
    log(`Gas used: ${receipt.gasUsed.toString()} units`);
    log(`Actual gas cost: ${ethers.formatEther(receipt.gasUsed * receipt.gasPrice)} ETH`);
    log(`Status: ${receipt.status === 1 ? 'SUCCESS' : 'FAILED'}`);
    
    if (receipt.status !== 1) {
      throw new Error('Transaction failed');
    }

    // Verify final balance
    const newWethBalance = await weth.balanceOf(wallet.address);
    log(`New WETH balance: ${ethers.formatEther(newWethBalance)} WETH`);
    log(`Transfer amount: ${ethers.formatEther(amountMinusGas)} WETH`);
    
    log('✅ Transfer completed successfully!');
    log(`Transaction: https://etherscan.io/tx/${tx.hash}`);
    
    process.exit(0);
    
  } catch (error) {
    log(`❌ Error: ${error.message}`);
    if (error.stack) {
      log(`Stack trace: ${error.stack}`);
    }
    process.exit(1);
  }
}

main();
