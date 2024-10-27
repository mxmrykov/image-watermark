package symbols

import "watermark/utils/ASCII"

type ASCII_SYM map[byte][][]bool

// GetASCIIRel - Getting lib of available symbols to write symbols
func GetASCIIRel() ASCII_SYM {
	newASCII := make(ASCII_SYM, len(ASCII.SymbolsRelations)+len(ASCII.SpecSymbols))
	// Add base uppercase english letters to our lib
	for i, sym := range ASCII.SymbolsRelations {
		newASCII[byte(65+i)] = sym
	}

	// Add special symbols, such as space
	for _, spec := range ASCII.SpecSymbols {
		newASCII[spec.ASCII_IX] = spec.RelationMatrix
	}

	return newASCII
}
