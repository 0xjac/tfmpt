
# TFMPT


The implementation is based on the following sources:

- The [Ethereum Yellow Paper specification, Appendix D.](https://ethereum.github.io/yellowpaper/paper.pdf)
- The [Ethereum Documentation](https://ethereum.org/en/developers/docs/data-structures-and-encoding/patricia-merkle-trie/)
- The [Go-Ethereum reference implementation](https://github.com/ethereum/go-ethereum/tree/master/trie)

## Disclaimer

After having spent a bit over one day implementing the trie,
I have realized that my intial estimate 
based on my theoretical knowledge of Merkle Patricia tries and Geth is innacurate.
I would need a few more days to implement the project fully.

I am happy to either discuss my implementation in its current state or spend a few more days. Note that the latter would require approval on my side, but I do not expect this to be an issue.

The following aspects are still missing:
    
- Deletion.
- Leaf nodes are not detected correctly during RLP decoding.
- Tests and docs.
- Computational Cost Analysis of Block Processing
