# Daily WETH Transfer Workflow

This workflow automatically transfers exactly 2 WETH from a configured wallet to `yaketh.eth` on Ethereum mainnet every day at 9:30 PM UTC. Gas fees are paid separately from the wallet's ETH balance.

## Overview

- **Start Date**: April 8, 2026 at 9:30 PM UTC
- **Schedule**: Daily at 21:30 UTC (9:30 PM)
- **Amount**: 2 WETH (gas paid separately from ETH balance)
- **Recipient**: yaketh.eth (ENS resolved)
- **Network**: Ethereum Mainnet

## Important Note

⚠️ **Current Date**: The system time shows April 9, 2026. Since the intended start date of April 8, 2026 at 9:30 PM UTC has passed, the workflow will begin executing from the next scheduled run (April 9, 2026 at 9:30 PM UTC) and continue daily thereafter.

## Setup Instructions

### 1. Configure GitHub Secrets

You need to add the following secrets to your repository:

Go to **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

#### Required Secrets:

1. **`TRANSFER_WALLET_PRIVATE_KEY`**
   - Description: Private key of the wallet that will send WETH
   - Format: Hex string (with or without `0x` prefix)
   - Example: `0x1234567890abcdef...`
   - ⚠️ **IMPORTANT**: Keep this secret secure! Never commit it to the repository.

2. **`ETHEREUM_RPC_URL`**
   - Description: Ethereum mainnet RPC endpoint URL
   - Format: HTTPS URL
   - Examples:
     - Infura: `https://mainnet.infura.io/v3/YOUR_PROJECT_ID`
     - Alchemy: `https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY`
     - Other providers: Any Ethereum mainnet RPC endpoint

### 2. Wallet Requirements

The wallet specified by `TRANSFER_WALLET_PRIVATE_KEY` must have:

1. **WETH Balance**: At least 2 WETH (the exact amount to transfer)
   - WETH9 contract: `0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2`
   - You can get WETH by wrapping ETH at https://weth.io or other DEXes

2. **ETH Balance**: Enough ETH to pay for gas fees (separate from WETH)
   - Recommended: At least 0.01 ETH for gas
   - Gas costs vary based on network congestion
   - Gas is paid in ETH, not deducted from the WETH transfer amount

### 3. Manual Triggering

The workflow can be manually triggered from the Actions tab:

1. Go to **Actions** tab in GitHub
2. Select **Daily WETH Transfer** workflow
3. Click **Run workflow** button
4. Select the branch and click **Run workflow**

### 4. Monitoring

#### View Logs

- Go to **Actions** tab
- Click on the workflow run
- View the logs for each step

#### Download Transfer Logs

- Each run creates a detailed log file
- Download from **Artifacts** section of the workflow run
- Logs are retained for 30 days

## How It Works

1. **Schedule Trigger**: GitHub Actions runs the workflow daily at 21:30 UTC
2. **ENS Resolution**: Resolves `yaketh.eth` to its Ethereum address
3. **Balance Check**: Verifies sufficient WETH and ETH balances
4. **Gas Estimation**: Estimates gas cost for the transfer
5. **Execute Transfer**: Sends exactly 2 WETH to the recipient (gas paid separately from ETH balance)
6. **Confirmation**: Waits for transaction confirmation
7. **Logging**: Records all details in a log file

## Security Considerations

1. **Private Key Storage**: 
   - Never commit private keys to the repository
   - Use GitHub Secrets for secure storage
   - Consider using a dedicated wallet for automated transfers

2. **Amount Limits**:
   - The script transfers a fixed amount (2 ETH by default)
   - Ensure the wallet has sufficient balance

3. **Gas Cost Protection**:
   - The script estimates gas costs before transfer
   - Gas is paid separately in ETH (not deducted from WETH)
   - Ensures the wallet has sufficient ETH for gas

4. **Network Verification**:
   - Script verifies it's connected to Ethereum mainnet (chainId 1)
   - Prevents accidental transfers on test networks

## Customization

### Change Transfer Amount

Edit `.github/workflows/daily-weth-transfer.yml`:

```yaml
env:
  AMOUNT_ETH: "2"  # Change to your desired amount
```

### Change Schedule

Edit the cron expression in `.github/workflows/daily-weth-transfer.yml`:

```yaml
schedule:
  - cron: '30 21 * * *'  # Current: 9:30 PM UTC daily
  # Format: minute hour day month weekday
  # Examples:
  # - cron: '0 0 * * *'   # Midnight UTC daily
  # - cron: '0 12 * * 1'  # Noon UTC every Monday
```

### Change Recipient

Edit `.github/workflows/daily-weth-transfer.yml`:

```yaml
env:
  RECIPIENT_ADDRESS: "yaketh.eth"  # Change to your desired ENS name or address
```

## Troubleshooting

### Common Issues

1. **"Insufficient WETH balance"**
   - Ensure wallet has at least 2 WETH
   - Wrap ETH to WETH if needed

2. **"Insufficient ETH for gas"**
   - Add more ETH to the wallet for gas fees
   - Recommended: Keep at least 0.01 ETH

3. **"Failed to resolve ENS name"**
   - Verify the ENS name is correct
   - Use a direct Ethereum address as alternative

4. **"Transaction failed"**
   - Check Etherscan for transaction details
   - Verify network isn't congested
   - Ensure gas price is sufficient

### Manual Testing

You can test the script locally:

```bash
# Install dependencies
npm install ethers@^6.13.0

# Set environment variables
export PRIVATE_KEY="your_private_key"
export RPC_URL="your_rpc_url"
export RECIPIENT_ADDRESS="yaketh.eth"
export AMOUNT_ETH="2"

# Run the script
node .github/scripts/transfer-weth.js
```

## Files

- `.github/workflows/daily-weth-transfer.yml` - GitHub Actions workflow configuration
- `.github/scripts/transfer-weth.js` - Transfer script implementation
- `.github/scripts/README.md` - This documentation

## Support

For issues or questions:
1. Check the workflow logs in the Actions tab
2. Review the transfer log artifacts
3. Verify wallet balances on Etherscan
4. Check RPC provider status

## License

This workflow is part of the eth2-beaconchain-explorer project.
