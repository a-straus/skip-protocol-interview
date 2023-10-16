# NFT Rarity Fetcher

A utility to fetch and analyze rarity scores for an NFT collection.

## Overview

This project fetches information about tokens and their traits, and then calculates rarity scores for each token. The tokens with the top rarity scores are then output in CSV format.

## Rarity Calculation

To calculate the rarity of a token, we look at each trait it possesses and calculate the inverse of the number of tokens having the same trait multiplied by the total unique traits in that category.

### Pseudocode

```
rarity = 0
for each trait category:
rarity += 1 / (tokens_with_same_trait * total_unique_traits_in_category)
```

### Example

Given a token:

```json
{
  "hat": "green beret",
  "earring": "gold"
}
```

In a collection where:

- 20 unique hat types exist and 50 tokens have a "green beret".
- 2 unique earring types exist and only 1 token has a "gold" earring.

The rarity score for this token would be:

`(1 / (50 _ 20)) + (1 / (1 _ 2))`

## Features

- Multi-threaded fetching of token data and rarity calculation
- Calculation of rarity scores for tokens based on trait frequency
- Output the top N rarities in CSV format

## Modules

### Main

The main module drives the program flow:

1. Parse and validate command-line flags.
2. Fetch all token traits and attributes.
3. Calculate the rarity for each token.
4. Output the top N rarities in CSV.

### Fetcher

The fetcher module is responsible for:

- Fetching token attributes concurrently using goroutines.
- Updating a shared map of traits and their frequencies.

### Calculator

The calculator module takes care of:

- Defining a minimum heap data structure to keep track of top N rarities.
- Calculating the rarity score of each token based on its attributes and the frequency of each trait in the collection.
- Providing the top N rarities based on the scores calculated.

### Logger

The logger module provides:

- A basic logging utility that outputs logs with time stamps.
- Features color-coded logs for different log levels such as Info, Error, and Success.

### Output

The output module handles:

- Writing the top N rarities to CSV output.
- CSV output includes Rank, Token ID, and Rarity.

# Running the program

```bash
go run main.go -flag1=value1 -flag2=value2 ...
```

(Replace `-flag1=value1 -flag2=value2` with actual flags and values required to run the program)

You can see a list of flags by running:

```bash
go run main.go --help
```
