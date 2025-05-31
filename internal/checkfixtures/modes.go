package main

// Problems whose outputs are not unique need either a semantic checker or an
// explicit comparison mode. Keep modes here so a mismatch is not silently
// treated as a valid alternative.
type compareMode int

const (
	modeNormal compareMode = iota
	modeUnorderedLines
	modeUnorderedTokens
	modeRotationLines
	modeDropFirstUnordered
	modeSmoke
)

var problemModes = map[string]compareMode{
	"course_4_algorithms_on_strings/week_1_suffix_trees/1_trie":                                                                   modeUnorderedLines,
	"course_4_algorithms_on_strings/week_1_suffix_trees/4_suffix_tree":                                                            modeUnorderedLines,
	"course_4_algorithms_on_strings/week_4_constructing_suffix_arrays_and_suffix_trees/3_suffix_array_matching":                   modeUnorderedTokens,
	"course_6_genome_assembly/week_1_assembling_the_phi_X174_using_overlap_graphs/1_small_genome_assembler_error_free":            modeRotationLines,
	"course_6_genome_assembly/week_1_assembling_the_phi_X174_using_overlap_graphs/2_small_genome_assembler_error_prone":           modeRotationLines,
	"course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/2_eulerian_cycle":                             modeRotationLines,
	"course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/4_genome_assembler_from_its_kmer_composition": modeRotationLines,
	"course_6_genome_assembly/week_3_faces_real_sequencing_data/5_genome_error_prone_assembler":                                   modeRotationLines,
}

func modeFor(rel string) compareMode {
	if mode, ok := problemModes[rel]; ok {
		return mode
	}
	return modeNormal
}
