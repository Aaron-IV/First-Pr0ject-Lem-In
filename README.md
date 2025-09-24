# Lem-in

A Go implementation of the Lem-in project: find the optimal way to move ants through a graph from a start room to an end room.

---

## 1. Usage

    ```sh
    go run cmd/main.go <input_file.txt>
    ```

    <input_file.txt> should describe the number of ants, rooms, and links in the required format.

## 2. Screenshot Examples
Example Input
```bash
    9
    #rooms
    ##start
    start 0 3
    ##end
    end 10 1
    C0 1 0
    C1 2 0
    C2 3 0
    C3 4 0
    I4 5 0
    I5 6 0
    A0 1 2
    A1 2 1
    A2 4 1
    B0 1 4
    B1 2 4
    E2 6 4
    D1 6 3
    D2 7 3
    D3 8 3
    H4 4 2
    H3 5 2
    F2 6 2
    F3 7 2
    F4 8 2
    G0 1 5
    G1 2 5
    G2 3 5
    G3 4 5
    G4 6 5
    H3-F2
    H3-H4
    H4-A2
    start-G0
    G0-G1
    G1-G2
    G2-G3
    G3-G4
    G4-D3
    start-A0
    A0-A1
    A0-D1
    A1-A2
    A1-B1
    A2-end
    A2-C3
    start-B0
    B0-B1
    B1-E2
    start-C0
    C0-C1
    C1-C2
    C2-C3
    C3-I4
    D1-D2
    D1-F2
    D2-E2
    D2-D3
    D2-F3
    D3-end
    F2-F3
    F3-F4
    F4-end
    I4-I5
    I5-end
```
Example Output
    ###############  4  ants #################
```bash
    L1->A0 L2->B0 L3->C0 
    L1->D1 L2->B1 L3->C1 L4->A0 L5->B0 
    L1->D2 L2->A1 L3->C2 L4->D1 L5->B1 L6->A0 L7->B0 
    L1->D3 L2->A2 L3->C3 L4->D2 L5->A1 L6->D1 L7->B1 L8->A0 L9->B0 
    L1->end L2->end L3->I4 L4->D3 L5->A2 L6->D2 L7->A1 L8->D1 L9->B1 
    L3->I5 L4->end L5->end L6->D3 L7->A2 L8->D2 L9->A1 
    L3->end L6->end L7->end L8->D3 L9->A2 
    L8->end L9->end 
```
## 3. Flags
    Currently, the program does not support command-line flags.
    You can add flags for visualization or debug output as needed.

## 4. Files 
    cmd/main.go — Main program entry point.
    helpers/ — Helper functions and types (e.g., parsing, utilities).
    visual/ — Visualization logic (if implemented).
    go.mod / go.sum — Go module files.
    README.md — This documentation file.
    input_file.txt — Example input file describing the graph.