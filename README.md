<!-- Copyright (C) 2024 Jacques Dafflon | 0xjac - All Rights Reserved -->

# TFMPT

The implementation is based on the following sources:

- The [Ethereum Yellow Paper specification, Appendix D.](https://ethereum.github.io/yellowpaper/paper.pdf)
- The [Ethereum Documentation](https://ethereum.org/en/developers/docs/data-structures-and-encoding/patricia-merkle-trie/)
- The [Go-Ethereum reference implementation](https://github.com/ethereum/go-ethereum/tree/master/trie)

> The implementation is not meant for production.
> It aims to demonstrate my ability to understand research papers and other sources,
> in order to correctly implement a data structure and related algorithms.

## Persistence

The persistence layer is implemented via [LevelDB](https://github.com/syndtr/goleveldb).
Alternative storage services can be easily used,
as long as they satisfy the simple [`store.DB` interface](./store/store.go).

Any setup or cleanup action required by the persistence layer must be handled outside this library.  

## Tests

The project uses [Task](https://taskfile.dev/) to run tests and coverage:

```console
go-task
task: Available tasks for this project:
* clean:       Clean up generated files.
* cover:       Run all the go tests with coverage.
* test:        Run all the go tests.
```

## Computational Cost Analysis of Block Processing

> During block processing, the Read/Write operation on state storage via disk are
> among the most computationally expensive operations.
>
> Define the following notation:
> - Let $R$ be the constant time to read one random element from persistent storage.
> - Let $W$ be the time to write an element.
> - Let $D_h$ be the total amount of state data for block $h \in N$.
>
> Further, assume:
> - Each block h consumes a constant amount of gas.
> - The distribution of opcodes across every smart contract of every block is the same.
> - The block period p (the minimum time between two mined blocks) is constant for all blocks h.
>
> Accounting for the MPT structure of the Ethereum state trie, derive $T_h$, the
> time required to process block $h$. Can you find a relation depending on $T_{h−1}$?  
> Remember to account for intermediate node lookups.

The time required to process block $h: T_h$ can be broken down into the time spent on reads and writes:

$$T_h=(\textrm{Time for reads})+((\textrm{Time for writes})$$

Considering lookups for a final (leaf) node results in lookups of intermediary nodes from the root to leaf, we define $d_r, d_w$ the average depth of nodes to traverse to read or write to a leaf. Considering the trie datastructure with $n$ nodes, then:
$$d_r \leq O(\log(n)), d_w \leq O(\log(n))$$  

Note that assuming homogenous data and no pre-existing data loaded in memory, $d_r \approx d_w$.

We define the number of read and write operations for block $h$ as $R_h$ and $W_h$ respectively.

Then $T_h$ is the number of read and write operations multiplied respectively by the average depth to traverse to read and write, multiplied by the time taken to read or write an element:

$$T_h = R_h \cdot d_r \cdot R + W_h \cdot d_w \cdot W$$

For $T_{h-1}$, considering each block consumes a constant amount of gas,
we assume the number of read and write operations to be similar for each block. Thus:
$$R_{h-1} \approx R_h \land W_{h-1} \approx W_h$$

We define $n_h$ the total number of nodes in the trie at block $h$
(for the state trie, this corresponds to the number of addresses used).
Assuming the number of insertions and deletions are constant, then: $n_h \approx n_{h-1}$. Hence, $d_r$ and $d_w$ are similar across blocks.

From our formula: $T_h = R_h \cdot d_r \cdot R + W_h \cdot d_w \cdot W$,
we have:
- $d_r$ and $d_w$ (roughly constant accross blocks),
- $R_{h-1} \approx R_h \land W_{h-1} \approx W_h$
- $R$ and $W$ are assumed constant.

Thus the resulting terms are roughly constant from block $h_1$ to block $h$.
That is:

$$T_h = R_h \cdot d_r \cdot R + W_h \cdot d_w \cdot W\hspace{2em}$$
$$\hspace{4em}\approx R_{h_1} \cdot d_r \cdot R + W_{h-1} \cdot d_w \cdot W = T_{h-1}$$

Hence $T_h \approx T_{h_1}$, the time to process a block is similar to, and *independent of* the time to process the previous block.
