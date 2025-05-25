# Problem index

One row per problem directory. `Tests` is the number of committed `test/N` fixtures. `PDF` links to the physical PDF page, not the printed page number. Read the PDF section for exact constraints and samples; the I/O column is only a shape reminder.

`Check` describes how `internal/checkfixtures` validates the output:

| Label | Meaning |
| --- | --- |
| `normal` | Whitespace-insensitive token compare. Output is unique. |
| `float` | Token compare with relative tolerance; default `-float 1e-6`. Pass a looser `-float` when the statement allows it. |
| `unordered lines` | Compare line sets; edge order does not matter. |
| `unordered tokens` | Compare token sets; position order does not matter. |
| `rotation` | Accept cyclic rotations of each line. Not arbitrary valid cycles, and not reverse complements. |
| `semantic` | Problem-specific invariant checker in `internal/checkfixtures/checker.go`. |
| `smoke` | Run only; exit status is checked, output is not compared. |
| `manual` | Output is not unique and no checker exists. The runner compares to one Go answer; a correct alternative can fail. Verify by hand or add a checker. |
| `normal + NGA50` | Three small fixtures use normal comparison; the visible dataset uses coverage and NGA50 against a reference. |

Fixture files are `test/N` plus `test/N.a`. Some are generated; see `docs/generated-fixtures.md`. Generated inputs are best-effort and can violate the original constraints. Go is the reference implementation; the Python solutions are legacy and `make fixtures-python` still fails on some non-unique and generated cases.

## Course 1 — Algorithmic Toolbox

### Week 1 — Programming challenges

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_sum_of_two_digits](../course_1_algorithmic_toolbox/week_1_programming_challenges/1_sum_of_two_digits) | 3 | [§1.1 · p. 8](../resources/course_1_assignments.pdf#page=8) | a b → a+b | `normal` |
| [2_maximum_pairwise_product](../course_1_algorithmic_toolbox/week_1_programming_challenges/2_maximum_pairwise_product) | 2 | [§1.2 · p. 12](../resources/course_1_assignments.pdf#page=12) | n, a[] → max a_i·a_j | `normal` |

### Week 2 — Algorithmic warm-up

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_fibonacci_number](../course_1_algorithmic_toolbox/week_2_algorithmic_warmup/1_fibonacci_number) | 6 | [§2.1.1 · p. 33](../resources/course_1_assignments.pdf#page=33) | n → F_n | `normal` |
| [2_last_digit_of_fibonacci_number](../course_1_algorithmic_toolbox/week_2_algorithmic_warmup/2_last_digit_of_fibonacci_number) | 3 | [§2.1.2 · p. 35](../resources/course_1_assignments.pdf#page=35) | n → F_n mod 10 | `normal` |
| [3_fibonacci_number_again](../course_1_algorithmic_toolbox/week_2_algorithmic_warmup/3_fibonacci_number_again) | 4 | [§2.1.3 · p. 37](../resources/course_1_assignments.pdf#page=37) | n m → F_n mod m | `normal` |
| [4_last_digit_of_the_sum_of_fibonacci_numbers](../course_1_algorithmic_toolbox/week_2_algorithmic_warmup/4_last_digit_of_the_sum_of_fibonacci_numbers) | 2 | [§2.1.4 · p. 39](../resources/course_1_assignments.pdf#page=39) | n → (Σ F_i) mod 10 | `normal` |
| [5_last_digit_of_the_sum_of_fibonacci_numbers_again](../course_1_algorithmic_toolbox/week_2_algorithmic_warmup/5_last_digit_of_the_sum_of_fibonacci_numbers_again) | 2 | [§2.1.5 · p. 40](../resources/course_1_assignments.pdf#page=40) | m n → (Σ F_m..F_n) mod 10 | `normal` |
| [6_last_digit_of_the_sum_of_squares_of_fibonacci_numbers](../course_1_algorithmic_toolbox/week_2_algorithmic_warmup/6_last_digit_of_the_sum_of_squares_of_fibonacci_numbers) | 3 | [§2.1.6 · p. 41](../resources/course_1_assignments.pdf#page=41) | n → (Σ F_i²) mod 10 | `normal` |
| [7_greatest_common_divisor](../course_1_algorithmic_toolbox/week_2_algorithmic_warmup/7_greatest_common_divisor) | 1 | [§2.1.7 · p. 43](../resources/course_1_assignments.pdf#page=43) | a b → gcd(a,b) | `normal` |
| [8_least_common_multiple](../course_1_algorithmic_toolbox/week_2_algorithmic_warmup/8_least_common_multiple) | 1 | [§2.1.8 · p. 44](../resources/course_1_assignments.pdf#page=44) | a b → lcm(a,b) | `normal` |

### Week 3 — Greedy algorithms

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_money_change](../course_1_algorithmic_toolbox/week_3_greedy_algorithms/1_money_change) | 2 | [§3.2.1 · p. 56](../resources/course_1_assignments.pdf#page=56) | money → min #coins 1,5,10 | `normal` |
| [2_maximum_value_of_the_loot](../course_1_algorithmic_toolbox/week_3_greedy_algorithms/2_maximum_value_of_the_loot) | 2 | [§3.2.2 · p. 57](../resources/course_1_assignments.pdf#page=57) | n W; (c_i,w_i)[] → max value | `float (stmt 1e-3)` |
| [3_car_fueling](../course_1_algorithmic_toolbox/week_3_greedy_algorithms/3_car_fueling) | 3 | [§3.2.3 · p. 59](../resources/course_1_assignments.pdf#page=59) | d, m, stops[] → min refills or -1 | `normal` |
| [4_maximum_advertisement_revenue](../course_1_algorithmic_toolbox/week_3_greedy_algorithms/4_maximum_advertisement_revenue) | 2 | [§3.2.4 · p. 61](../resources/course_1_assignments.pdf#page=61) | n, price[], clicks[] → max Σ price·click | `normal` |
| [5_collecting_signatures](../course_1_algorithmic_toolbox/week_3_greedy_algorithms/5_collecting_signatures) | 2 | [§3.2.5 · p. 63](../resources/course_1_assignments.pdf#page=63) | n, [l_i,r_i][] → k; k points | `manual` |
| [6_maximum_number_of_prizes](../course_1_algorithmic_toolbox/week_3_greedy_algorithms/6_maximum_number_of_prizes) | 3 | [§3.2.6 · p. 65](../resources/course_1_assignments.pdf#page=65) | n → k; k distinct ints summing n | `manual` |
| [7_maximum_salary](../course_1_algorithmic_toolbox/week_3_greedy_algorithms/7_maximum_salary) | 3 | [§3.2.7 · p. 67](../resources/course_1_assignments.pdf#page=67) | n, a[] → largest concatenation | `normal` |

### Week 4 — Divide and conquer

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_binary_search](../course_1_algorithmic_toolbox/week_4_divide_and_conquer/1_binary_search) | 1 | [§4.2.1 · p. 99](../resources/course_1_assignments.pdf#page=99) | n, k[], m, q[] → index or -1 | `normal` |
| [2_binary_search_duplicates](../course_1_algorithmic_toolbox/week_4_divide_and_conquer/2_binary_search_duplicates) | 1 | [§4.2.2 · p. 102](../resources/course_1_assignments.pdf#page=102) | n, k[], m, q[] → first index or -1 | `normal` |
| [3_majority_element](../course_1_algorithmic_toolbox/week_4_divide_and_conquer/3_majority_element) | 3 | [§4.2.3 · p. 104](../resources/course_1_assignments.pdf#page=104) | n, a[] → 1 if majority else 0 | `normal` |
| [4_improving_quicksort](../course_1_algorithmic_toolbox/week_4_divide_and_conquer/4_improving_quicksort) | 2 | [§4.2.4 · p. 105](../resources/course_1_assignments.pdf#page=105) | n, a[] → sorted a[] | `normal` |
| [5_number_of_inversions](../course_1_algorithmic_toolbox/week_4_divide_and_conquer/5_number_of_inversions) | 1 | [§4.2.5 · p. 107](../resources/course_1_assignments.pdf#page=107) | n, a[] → inversion count | `normal` |
| [6_organizing_a_lottery](../course_1_algorithmic_toolbox/week_4_divide_and_conquer/6_organizing_a_lottery) | 3 | [§4.2.6 · p. 109](../resources/course_1_assignments.pdf#page=109) | n m; segments; points → counts | `normal` |
| [7_closest_points](../course_1_algorithmic_toolbox/week_4_divide_and_conquer/7_closest_points) | 2 | [§4.2.7 · p. 111](../resources/course_1_assignments.pdf#page=111) | n; points → min distance | `float (stmt 1e-3)` |

### Week 5 — Dynamic programming 1

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_money_change_again](../course_1_algorithmic_toolbox/week_5_dynamic_programming_1/1_money_change_again) | 1 | [§5.2.1 · p. 131](../resources/course_1_assignments.pdf#page=131) | money → min #coins 1,3,4 | `normal` |
| [2_primitive_calculator](../course_1_algorithmic_toolbox/week_5_dynamic_programming_1/2_primitive_calculator) | 2 | [§5.2.2 · p. 132](../resources/course_1_assignments.pdf#page=132) | n → k; sequence 1→n | `manual` |
| [3_edit_distance](../course_1_algorithmic_toolbox/week_5_dynamic_programming_1/3_edit_distance) | 3 | [§5.2.3 · p. 134](../resources/course_1_assignments.pdf#page=134) | s t → edit distance | `normal` |
| [4_longest_common_subsequence_of_two_sequences](../course_1_algorithmic_toolbox/week_5_dynamic_programming_1/4_longest_common_subsequence_of_two_sequences) | 3 | [§5.2.4 · p. 136](../resources/course_1_assignments.pdf#page=136) | n, a[]; m, b[] → LCS length | `normal` |
| [5_longest_common_subsequence_of_three_sequences](../course_1_algorithmic_toolbox/week_5_dynamic_programming_1/5_longest_common_subsequence_of_three_sequences) | 2 | [§5.2.5 · p. 138](../resources/course_1_assignments.pdf#page=138) | n, a[]; m, b[]; l, c[] → LCS length | `normal` |

### Week 6 — Dynamic programming 2

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_maximum_amount_of_gold](../course_1_algorithmic_toolbox/week_6_dynamic_programming_2/1_maximum_amount_of_gold) | 1 | [§5.2.6 · p. 140](../resources/course_1_assignments.pdf#page=140) | W n; weights[] → max weight ≤ W | `normal` |
| [2_partitioning_souvenirs](../course_1_algorithmic_toolbox/week_6_dynamic_programming_2/2_partitioning_souvenirs) | 4 | [§5.2.7 · p. 142](../resources/course_1_assignments.pdf#page=142) | n, v[] → 1 if three equal sums else 0 | `normal` |
| [3_maximum_value_of_an_arithmetic_expression](../course_1_algorithmic_toolbox/week_6_dynamic_programming_2/3_maximum_value_of_an_arithmetic_expression) | 1 | [§5.2.8 · p. 144](../resources/course_1_assignments.pdf#page=144) | expression → max value | `normal` |

## Course 2 — Data Structures

### Week 1 — Basic data structures

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_brackets_in_code](../course_2_data_structures/week_1_basic_data_structures/1_brackets_in_code) | 51 | [§1 · p. 3](../resources/course_2_assignments_week_1.pdf#page=3) | string S → Success or 1-based error index | `normal` |
| [2_tree_height](../course_2_data_structures/week_1_basic_data_structures/2_tree_height) | 16 | [§2 · p. 6](../resources/course_2_assignments_week_1.pdf#page=6) | n; parent[] → height | `normal` |
| [3_network_simulation](../course_2_data_structures/week_1_basic_data_structures/3_network_simulation) | 18 | [§3 · p. 9](../resources/course_2_assignments_week_1.pdf#page=9) | S n; n packets (A_i,P_i) → start times/-1 | `normal` |
| [4_stack_with_max](../course_2_data_structures/week_1_basic_data_structures/4_stack_with_max) | 5 | [§4 · p. 12](../resources/course_2_assignments_week_1.pdf#page=12) | q; push/pop/max → max answers | `normal` |
| [5_max_sliding_window](../course_2_data_structures/week_1_basic_data_structures/5_max_sliding_window) | 1 | [§5 · p. 15](../resources/course_2_assignments_week_1.pdf#page=15) | n, a[]; m → window maxima | `normal` |

### Week 2 — Priority queues and disjoint sets

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_make_heap](../course_2_data_structures/week_2_priority_queues_and_disjoint_sets/1_make_heap) | 2 | [§1 · p. 2](../resources/course_2_assignments_week_2.pdf#page=2) | n, a[] → m; m swap pairs | `manual` |
| [2_job_queue](../course_2_data_structures/week_2_priority_queues_and_disjoint_sets/2_job_queue) | 2 | [§2 · p. 4](../resources/course_2_assignments_week_2.pdf#page=4) | n m; t[] → thread,start per job | `normal` |
| [3_merging_tables](../course_2_data_structures/week_2_priority_queues_and_disjoint_sets/3_merging_tables) | 2 | [§3 · p. 6](../resources/course_2_assignments_week_2.pdf#page=6) | n m; r[]; merge queries → max sizes | `normal` |

### Week 3 — Hash tables

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_phone_book](../course_2_data_structures/week_3_hash_tables/1_phone_book) | 2 | [§1 · p. 2](../resources/course_2_assignments_week_3.pdf#page=2) | N; add/find/del → names/not found | `normal` |
| [2_hash_chains](../course_2_data_structures/week_3_hash_tables/2_hash_chains) | 4 | [§2 · p. 5](../resources/course_2_assignments_week_3.pdf#page=5) | m; N; add/del/find/check → chain answers | `normal` |
| [3_hash_substring](../course_2_data_structures/week_3_hash_tables/3_hash_substring) | 4 | [§3 · p. 9](../resources/course_2_assignments_week_3.pdf#page=9) | P T → occurrence positions | `normal` |
| [4_substring_equality](../course_2_data_structures/week_3_hash_tables/4_substring_equality) | 1 | [§4 · p. 11](../resources/course_2_assignments_week_3.pdf#page=11) | s; q; (a,b,l)[] → Yes/No | `normal` |
| [5_longest_common_substring](../course_2_data_structures/week_3_hash_tables/5_longest_common_substring) | 2 | [§5 · p. 14](../resources/course_2_assignments_week_3.pdf#page=14) | s t per line → i j len | `semantic` |
| [6_matching_with_mismatches](../course_2_data_structures/week_3_hash_tables/6_matching_with_mismatches) | 3 | [§6 · p. 16](../resources/course_2_assignments_week_3.pdf#page=16) | k t p per line → count; positions | `normal` |

### Week 4 — Binary search trees

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_tree_traversals](../course_2_data_structures/week_4_binary_search_trees/1_tree_traversals) | 2 | [§1 · p. 2](../resources/course_2_assignments_week_4.pdf#page=2) | n; (key,left,right)[] → in/pre/post | `normal` |
| [2_is_bst](../course_2_data_structures/week_4_binary_search_trees/2_is_bst) | 7 | [§2 · p. 5](../resources/course_2_assignments_week_4.pdf#page=5) | n; nodes[] → CORRECT/INCORRECT | `normal` |
| [3_is_bst_advanced](../course_2_data_structures/week_4_binary_search_trees/3_is_bst_advanced) | 9 | [§3 · p. 9](../resources/course_2_assignments_week_4.pdf#page=9) | n; nodes[] → CORRECT/INCORRECT | `normal` |
| [4_set_range_sum](../course_2_data_structures/week_4_binary_search_trees/4_set_range_sum) | 6 | [§4 · p. 14](../resources/course_2_assignments_week_4.pdf#page=14) | n; +/-/?/s ops → Found/Not found/sum | `normal` |
| [5_rope](../course_2_data_structures/week_4_binary_search_trees/5_rope) | 2 | [§5 · p. 18](../resources/course_2_assignments_week_4.pdf#page=18) | S; q; (i,j,k)[] → final string | `normal` |

## Course 3 — Algorithms on Graphs

### Week 1 — Graph decomposition 1

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_reachability](../course_3_algorithms_on_graphs/week_1_graph_decomposition1/1_reachability) | 2 | [§1 · p. 5](../resources/course_3_assignments_week_1.pdf#page=5) | n m; edges; u v → 1/0 | `normal` |
| [2_connected_components](../course_3_algorithms_on_graphs/week_1_graph_decomposition1/2_connected_components) | 1 | [§2 · p. 7](../resources/course_3_assignments_week_1.pdf#page=7) | n m; edges → component count | `normal` |

### Week 2 — Graph decomposition 2

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_cs_curriculum](../course_3_algorithms_on_graphs/week_2_graph_decomposition2/1_cs_curriculum) | 2 | [§1 · p. 5](../resources/course_3_assignments_week_2.pdf#page=5) | n m; edges → 1 if cycle else 0 | `normal` |
| [2_order_of_courses](../course_3_algorithms_on_graphs/week_2_graph_decomposition2/2_order_of_courses) | 3 | [§2 · p. 7](../resources/course_3_assignments_week_2.pdf#page=7) | n m; DAG edges → topological order | `semantic` |
| [3_intersection_reachability](../course_3_algorithms_on_graphs/week_2_graph_decomposition2/3_intersection_reachability) | 2 | [§3 · p. 9](../resources/course_3_assignments_week_2.pdf#page=9) | n m; edges → n×n reachability matrix | `normal` |

### Week 3 — Paths in graphs 1

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_flight_segments](../course_3_algorithms_on_graphs/week_3_paths_in_graphs1/1_flight_segments) | 2 | [§1 · p. 4](../resources/course_3_assignments_week_3.pdf#page=4) | n m; edges; s t → min edges or -1 | `normal` |
| [2_bipartite](../course_3_algorithms_on_graphs/week_3_paths_in_graphs1/2_bipartite) | 2 | [§2 · p. 6](../resources/course_3_assignments_week_3.pdf#page=6) | n m; edges → 1/0 | `normal` |

### Week 4 — Paths in graphs 2

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_minimum_flight_cost](../course_3_algorithms_on_graphs/week_4_paths_in_graphs2/1_minimum_flight_cost) | 3 | [§1 · p. 4](../resources/course_3_assignments_week_4.pdf#page=4) | n m; edges; s t → min cost or -1 | `normal` |
| [2_detecting_anomalies](../course_3_algorithms_on_graphs/week_4_paths_in_graphs2/2_detecting_anomalies) | 2 | [§2 · p. 6](../resources/course_3_assignments_week_4.pdf#page=6) | n m; edges → 1 if negative cycle else 0 | `normal` |
| [3_exchanging_money](../course_3_algorithms_on_graphs/week_4_paths_in_graphs2/3_exchanging_money) | 3 | [§3 · p. 7](../resources/course_3_assignments_week_4.pdf#page=7) | n m; edges; s → distances (-, *) | `normal` |

### Week 5 — Spanning trees

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_connecting_points](../course_3_algorithms_on_graphs/week_5_spanning_trees/1_connecting_points) | 2 | [§1 · p. 4](../resources/course_3_assignments_week_5.pdf#page=4) | n; points → MST cost | `float` |
| [2_clustering](../course_3_algorithms_on_graphs/week_5_spanning_trees/2_clustering) | 2 | [§2 · p. 6](../resources/course_3_assignments_week_5.pdf#page=6) | n; points; k → max spacing | `float` |

### Week 6 — Advanced shortest paths

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_friend_suggestion](../course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/1_friend_suggestion) | 2 | [§1 · p. 3](../resources/course_3_assignments_week_6.pdf#page=3) | n m; edges; q queries → distance/count | `normal` |
| [2_dist_with_coords](../course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/2_dist_with_coords) | 3 | [§2 · p. 6](../resources/course_3_assignments_week_6.pdf#page=6) | n m; coords; edges; q → distances | `normal` |
| [3_dist_preprocess_small](../course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/3_dist_preprocess_small) | 8 | [§3 · p. 9](../resources/course_3_assignments_week_6.pdf#page=9) | n m; edges; q → distances | `normal` |
| [4_dist_preprocess_large](../course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/4_dist_preprocess_large) | 10 | [§4 · p. 12](../resources/course_3_assignments_week_6.pdf#page=12) | n m; edges; q → distances | `normal` |
| [5_travelling_salesman_problem](../course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/5_travelling_salesman_problem) | 9 | [§5 · p. 14](../resources/course_3_assignments_week_6.pdf#page=14) | n m; edges → min cost; tour | `manual` |

## Course 4 — Algorithms on Strings

### Week 1 — Suffix trees

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_trie](../course_4_algorithms_on_strings/week_1_suffix_trees/1_trie) | 3 | [§1 · p. 3](../resources/course_4_assignments_week_1.pdf#page=3) | patterns[] → edges u->v:c | `unordered lines` |
| [2_trie_matching](../course_4_algorithms_on_strings/week_1_suffix_trees/2_trie_matching) | 3 | [§2 · p. 6](../resources/course_4_assignments_week_1.pdf#page=6) | Text; n; patterns[] → positions | `normal` |
| [3_trie_matching_extended](../course_4_algorithms_on_strings/week_1_suffix_trees/3_trie_matching_extended) | 2 | [§3 · p. 8](../resources/course_4_assignments_week_1.pdf#page=8) | Text; n; patterns[] → positions | `normal` |
| [4_suffix_tree](../course_4_algorithms_on_strings/week_1_suffix_trees/4_suffix_tree) | 3 | [§4 · p. 10](../resources/course_4_assignments_week_1.pdf#page=10) | Text$ → edge labels | `unordered lines` |
| [5_non_shared_substring](../course_4_algorithms_on_strings/week_1_suffix_trees/5_non_shared_substring) | 4 | [§5 · p. 13](../resources/course_4_assignments_week_1.pdf#page=13) | Text1 Text2 → shortest non-shared substring | `semantic` |

### Week 2 — Burrows–Wheeler transform and suffix arrays

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_bwt](../course_4_algorithms_on_strings/week_2_burrows_wheeler_transform_and_suffix_arrays/1_bwt) | 3 | [§1 · p. 3](../resources/course_4_assignments_week_2.pdf#page=3) | Text$ → BWT(Text) | `normal` |
| [2_bwtinverse](../course_4_algorithms_on_strings/week_2_burrows_wheeler_transform_and_suffix_arrays/2_bwtinverse) | 2 | [§2 · p. 5](../resources/course_4_assignments_week_2.pdf#page=5) | Transform → Text | `normal` |
| [3_bwmatching](../course_4_algorithms_on_strings/week_2_burrows_wheeler_transform_and_suffix_arrays/3_bwmatching) | 4 | [§3 · p. 7](../resources/course_4_assignments_week_2.pdf#page=7) | BWT; n; patterns[] → match counts | `normal` |
| [4_suffix_array](../course_4_algorithms_on_strings/week_2_burrows_wheeler_transform_and_suffix_arrays/4_suffix_array) | 3 | [§4 · p. 10](../resources/course_4_assignments_week_2.pdf#page=10) | Text$ → suffix array | `normal` |

### Week 4 — Constructing suffix arrays and suffix trees

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_kmp](../course_4_algorithms_on_strings/week_4_constructing_suffix_arrays_and_suffix_trees/1_kmp) | 3 | [§1 · p. 3](../resources/course_4_assignments_week_3.pdf#page=3) | Pattern Genome → positions | `normal` |
| [2_suffix_array_long](../course_4_algorithms_on_strings/week_4_constructing_suffix_arrays_and_suffix_trees/2_suffix_array_long) | 4 | [§2 · p. 5](../resources/course_4_assignments_week_3.pdf#page=5) | Text$ → suffix array | `normal` |
| [3_suffix_array_matching](../course_4_algorithms_on_strings/week_4_constructing_suffix_arrays_and_suffix_trees/3_suffix_array_matching) | 3 | [§3 · p. 7](../resources/course_4_assignments_week_3.pdf#page=7) | Text; n; patterns[] → positions | `unordered tokens` |
| [4_suffix_tree_from_array](../course_4_algorithms_on_strings/week_4_constructing_suffix_arrays_and_suffix_trees/4_suffix_tree_from_array) | 3 | [§4 · p. 9](../resources/course_4_assignments_week_3.pdf#page=9) | Text$; SA; LCP → Text + edges | `smoke` |

## Course 5 — Advanced Algorithms and Complexity

### Week 1 — Flows in networks

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_evacuation](../course_5_advanced_algorithms_and_complexity/week_1_flows_in_networks/1_evacuation) | 19 | [§1 · p. 3](../resources/course_5_assignments_week_1.pdf#page=3) | n m; edges (u,v,c) → max flow | `normal` |
| [2_airline_crews](../course_5_advanced_algorithms_and_complexity/week_1_flows_in_networks/2_airline_crews) | 20 | [§2 · p. 6](../resources/course_5_assignments_week_1.pdf#page=6) | n m; matrix → crew per flight | `manual` |
| [3_stock_charts](../course_5_advanced_algorithms_and_complexity/week_1_flows_in_networks/3_stock_charts) | 36 | [§3 · p. 8](../resources/course_5_assignments_week_1.pdf#page=8) | n k; prices[] → min charts | `normal` |

### Week 2 — Linear programming

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_energy_values](../course_5_advanced_algorithms_and_complexity/week_2_linear_programming/1_energy_values) | 4 | [§1 · p. 3](../resources/course_5_assignments_week_2.pdf#page=3) | n; n×n matrix → values | `float (stmt 1e-2)` |
| [2_diet](../course_5_advanced_algorithms_and_complexity/week_2_linear_programming/2_diet) | 8 | [§2 · p. 6](../resources/course_5_assignments_week_2.pdf#page=6) | n m; A,b,c → status + x[] | `manual` |
| [3_ad_allocation](../course_5_advanced_algorithms_and_complexity/week_2_linear_programming/3_ad_allocation) | 9 | [§3 · p. 10](../resources/course_5_assignments_week_2.pdf#page=10) | p q; A,b,c → status + x[] | `manual` |

### Week 3 — NP-complete problems

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_gsm_network](../course_5_advanced_algorithms_and_complexity/week_3_NP_complete_problems/1_gsm_network) | 2 | [§1 · p. 3](../resources/course_5_assignments_week_3.pdf#page=3) | n m; edges → CNF | `smoke` |
| [2_cleaning_apartment](../course_5_advanced_algorithms_and_complexity/week_3_NP_complete_problems/2_cleaning_apartment) | 2 | [§2 · p. 6](../resources/course_5_assignments_week_3.pdf#page=6) | n m; edges → CNF | `smoke` |
| [3_budget_allocation](../course_5_advanced_algorithms_and_complexity/week_3_NP_complete_problems/3_budget_allocation) | 2 | [§3 · p. 9](../resources/course_5_assignments_week_3.pdf#page=9) | n m; A,b → CNF | `manual` |

### Week 4 — Coping with NP-completeness

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_circuit_design](../course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/1_circuit_design) | 2 | [§1 · p. 3](../resources/course_5_assignments_week_4.pdf#page=3) | V C; clauses → SAT/UNSAT + assignment | `manual` |
| [2_plan_party](../course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/2_plan_party) | 3 | [§2 · p. 5](../resources/course_5_assignments_week_4.pdf#page=5) | n; f[]; tree edges → max fun | `normal` |
| [3_school_bus](../course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/3_school_bus) | 2 | [§3 · p. 8](../resources/course_5_assignments_week_4.pdf#page=8) | n m; edges → min cost; tour | `manual` |
| [4_reschedule_exams](../course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/4_reschedule_exams) | 4 | [§4 · p. 11](../resources/course_5_assignments_week_4.pdf#page=11) | n m; colors; edges → Impossible or colors | `manual` |

### Week 5 — Streaming algorithms

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_heavy_hitters](../course_5_advanced_algorithms_and_complexity/week_5_streaming_algorithms/1_heavy_hitters) | 10 | [§1 · p. 2](../resources/course_5_assignments_week_5.pdf#page=2) | N t; A; B; q ids → bits | `normal` |

## Course 6 — Genome Assembly Programming Challenge

### Week 1 — Assembling the phi X174 using overlap graphs

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_small_genome_assembler_error_free](../course_6_genome_assembly/week_1_assembling_the_phi_X174_using_overlap_graphs/1_small_genome_assembler_error_free) | 1 | [§1 · p. 3](../resources/course_6_assignments_week_1.pdf#page=3) | reads[] → genome | `rotation` |
| [2_small_genome_assembler_error_prone](../course_6_genome_assembly/week_1_assembling_the_phi_X174_using_overlap_graphs/2_small_genome_assembler_error_prone) | 2 | [§2 · p. 5](../resources/course_6_assignments_week_1.pdf#page=5) | reads[] → genome | `rotation` |

### Week 2 — Assembling the phi X174 using de bruijn graphs

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_puzzle_assembly](../course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/1_puzzle_assembly) | 2 | [§1 · p. 3](../resources/course_6_assignments_week_2.pdf#page=3) | 25 pieces → 5×5 placement | `manual` |
| [2_eulerian_cycle](../course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/2_eulerian_cycle) | 5 | [§2 · p. 5](../resources/course_6_assignments_week_2.pdf#page=5) | n m; edges → 1 + cycle or 0 | `rotation` |
| [3_k_universal_circular_string](../course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/3_k_universal_circular_string) | 2 | [§3 · p. 8](../resources/course_6_assignments_week_2.pdf#page=8) | k → binary circular string | `semantic` |
| [4_genome_assembler_from_its_kmer_composition](../course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/4_genome_assembler_from_its_kmer_composition) | 2 | [§4 · p. 9](../resources/course_6_assignments_week_2.pdf#page=9) | k-mers[] → genome | `rotation` |

### Week 3 — Faces real sequencing data

| Problem | Tests | PDF | Input → output | Check |
| --- | ---: | --- | --- | --- |
| [1_network_circulation](../course_6_genome_assembly/week_3_faces_real_sequencing_data/1_network_circulation) | 3 | [§1 · p. 3](../resources/course_6_assignments_week_3.pdf#page=3) | n m; (u,v,l,c)[] → YES + flows or NO | `manual` |
| [2_optimal_kmer_size](../course_6_genome_assembly/week_3_faces_real_sequencing_data/2_optimal_kmer_size) | 1 | [§2 · p. 5](../resources/course_6_assignments_week_3.pdf#page=5) | reads[] → k | `manual` |
| [3_bubble_detection](../course_6_genome_assembly/week_3_faces_real_sequencing_data/3_bubble_detection) | 2 | [§3 · p. 7](../resources/course_6_assignments_week_3.pdf#page=7) | t; reads[] → bubble count | `normal` |
| [4_tip_removal](../course_6_genome_assembly/week_3_faces_real_sequencing_data/4_tip_removal) | 4 | [§4 · p. 9](../resources/course_6_assignments_week_3.pdf#page=9) | reads[] → removed-edge count | `normal` |
| [5_genome_error_prone_assembler](../course_6_genome_assembly/week_3_faces_real_sequencing_data/5_genome_error_prone_assembler) | 1 | [§5 · p. 11](../resources/course_6_assignments_week_3.pdf#page=11) | reads[] → genome | `rotation` |
| [6_implementing_an_assembler](../course_6_genome_assembly/week_3_faces_real_sequencing_data/6_implementing_an_assembler) | 3 | [§6 · p. 12](../resources/course_6_assignments_week_3.pdf#page=12) | reads or read-pairs → FASTA contigs | `normal + NGA50` |

## Non-unique outputs

These problems have more than one correct output. Semantic and order-insensitive checks validate content; rotation checks only cyclic shifts; smoke checks only exit status; `manual` rows still compare against one Go answer.

**Semantic checker (4):** `5_longest_common_substring`, `2_order_of_courses`, `5_non_shared_substring`, `3_k_universal_circular_string`.

**Order-insensitive (3):** `1_trie`, `4_suffix_tree`, `3_suffix_array_matching`.

**Rotation-only (5):** `1_small_genome_assembler_error_free`, `2_small_genome_assembler_error_prone`, `2_eulerian_cycle`, `4_genome_assembler_from_its_kmer_composition`, `5_genome_error_prone_assembler`.

**Smoke-only (3):** `4_suffix_tree_from_array`, `1_gsm_network`, `2_cleaning_apartment`.

**No checker, one-answer compare (15):**

- `course_1_algorithmic_toolbox/week_3_greedy_algorithms/5_collecting_signatures` — n, [l_i,r_i][] → k; k points
- `course_1_algorithmic_toolbox/week_3_greedy_algorithms/6_maximum_number_of_prizes` — n → k; k distinct ints summing n
- `course_1_algorithmic_toolbox/week_5_dynamic_programming_1/2_primitive_calculator` — n → k; sequence 1→n
- `course_2_data_structures/week_2_priority_queues_and_disjoint_sets/1_make_heap` — n, a[] → m; m swap pairs
- `course_3_algorithms_on_graphs/week_6_advanced_shortest_paths/5_travelling_salesman_problem` — n m; edges → min cost; tour
- `course_5_advanced_algorithms_and_complexity/week_1_flows_in_networks/2_airline_crews` — n m; matrix → crew per flight
- `course_5_advanced_algorithms_and_complexity/week_2_linear_programming/2_diet` — n m; A,b,c → status + x[]
- `course_5_advanced_algorithms_and_complexity/week_2_linear_programming/3_ad_allocation` — p q; A,b,c → status + x[]
- `course_5_advanced_algorithms_and_complexity/week_3_NP_complete_problems/3_budget_allocation` — n m; A,b → CNF
- `course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/1_circuit_design` — V C; clauses → SAT/UNSAT + assignment
- `course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/3_school_bus` — n m; edges → min cost; tour
- `course_5_advanced_algorithms_and_complexity/week_4_coping_with_NP_completeness/4_reschedule_exams` — n m; colors; edges → Impossible or colors
- `course_6_genome_assembly/week_2_assembling_the_phi_X174_using_de_bruijn_graphs/1_puzzle_assembly` — 25 pieces → 5×5 placement
- `course_6_genome_assembly/week_3_faces_real_sequencing_data/1_network_circulation` — n m; (u,v,l,c)[] → YES + flows or NO
- `course_6_genome_assembly/week_3_faces_real_sequencing_data/2_optimal_kmer_size` — reads[] → k

The final project is separate: the visible dataset uses a local coverage/NGA50 approximation of the Coursera grader, not an exact-output checker.
