compile:
	solc --abi contracts/Token.sol -o build --overwrite
	solc --bin contracts/Token.sol -o build --overwrite
	abigen --bin build/MyToken.bin --abi build/MyToken.abi --pkg token --out contracts/Token.go