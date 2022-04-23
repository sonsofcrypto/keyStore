# KeyStore

Stores private keys according to [Web3 Secret Storage Definition](https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition).
Extends format to optionally store mnemonic similarly to [Ethers.js](https://github.com/ethers-io/ethers.js/blob/f599d6f23dad0d0acaa3828d6b7acaab2d5e455b/packages/json-wallets/src.ts/keystore.ts#L364).
Supports storing of `KeyStoreItem`s to variety `Backend`s (eg disk, system keychain).

### TODO: 
    - [ ] Add caching to `KeyStore` operations
    - [ ] Get standalone module for secp256k1 instead of importing "github.com/ethereum/go-ethereum/crypto"


