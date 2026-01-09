package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "path/filepath"
    //"log"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println(" --------------------------------------------------------")
        fmt.Println(" CLEAN_WATERS deletes water molecules with zero occupancy")
        fmt.Println(" Usage: clean_waters <INPUT>")
        fmt.Println(" --------------------------------------------------------")
        return
    }

    // -------------------------
    // CHECK NUMBER OF ARGUMENTS

    inputFile := os.Args[1]
    outputFile := "cleaned_" + inputFile

    // ---------------------
    // TEST READ/WRITE FILES

    inFile, err := os.Open(inputFile)
    if err != nil {
        fmt.Println("Error opening input file:", err)
        return
    }
    defer inFile.Close()

    outFile, err := os.Create(outputFile)
    if err != nil {
        fmt.Println("Error creating output file:", err)
        return
    }
    defer outFile.Close()

    // -------------------------
    // ANALYZE INPUT FILES

    ext := filepath.Ext(inputFile)
    switch ext {
    
    // PDB
    case ".pdb":
        fmt.Println("Detected file type: PDB")
        scanner := bufio.NewScanner(inFile)
        writer := bufio.NewWriter(outFile)

        fmt.Println("Deleted lines (HOH with zero occupancy):")

        for scanner.Scan() {
            line := scanner.Text()

            if (strings.HasPrefix(line, "HETATM") || strings.HasPrefix(line, "ATOM") )  && strings.Contains(line[17:20], "HOH") {
                occupancyStr := strings.TrimSpace(line[54:60])
                occupancy, err := strconv.ParseFloat(occupancyStr, 64)
                if err != nil || occupancy == 0.0 {
                    fmt.Println(line) // Report the deleted line
                    continue          // Skip writing this line
                }
            }

            writer.WriteString(line + "\n")
        }

        if err := scanner.Err(); err != nil {
            fmt.Println("Error reading input file:", err)
        }

        writer.Flush()
        fmt.Printf("Cleaning complete. Output written to %s\n", outputFile)


    case ".cif":
        fmt.Println("Detected file type: CIF")
        scanner := bufio.NewScanner(inFile)
        writer := bufio.NewWriter(outFile)
        position := -1
        position_search := -1

        fmt.Println("Deleted lines (HOH with zero occupancy):")        
        
        for scanner.Scan() {
            line := scanner.Text()
            token := strings.Fields(line)
            if strings.HasPrefix(line, "loop_") {
                position_search = -1
            }
            if strings.HasPrefix(line, "_atom_site.occupancy") {
                position = position_search
            }
            
            if (strings.HasPrefix(line, "HETATM") || strings.HasPrefix(line, "ATOM")) && token[position] == "0" && strings.Contains(line, "HOH") {
                fmt.Println(line)
                continue
            }
            
            writer.WriteString(line + "\n")
            position_search++
        }
        
        
    case ".mmcif":
        fmt.Println("Detected file type: mmCIF")
        scanner := bufio.NewScanner(inFile)
        writer := bufio.NewWriter(outFile)
        position := -1
        position_search := -1

        fmt.Println("Deleted lines (HOH with zero occupancy):")        
        
        for scanner.Scan() {
            line := scanner.Text()
            token := strings.Fields(line)
            if strings.HasPrefix(line, "loop_") {
                position_search = -1
            }
            if strings.HasPrefix(line, "_atom_site.occupancy") {
                position = position_search
            }
            
            if (strings.HasPrefix(line, "HETATM") || strings.HasPrefix(line, "ATOM")) && token[position] == "0" && strings.Contains(line, "HOH") {
                fmt.Println(line)
                continue
            }
            
            writer.WriteString(line + "\n")
            position_search++
        }


    default:
        fmt.Printf("Unknown file type: %s\n", ext)
        return
    }


}
