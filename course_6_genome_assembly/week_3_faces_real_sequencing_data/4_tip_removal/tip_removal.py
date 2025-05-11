# python3

# given a list of error-prone reads, construct a de bruijn graph from the
# 15-mers created from the reads and perform the task of tip removal on this de
# bruijn graph


import sys


class TipRemoval:
    def __init__(self, reads, k):
        self.k = k
        self.threshold = k
        self.graph = {}
        self.edges_removed = 0
        self.build_de_bruijn_graph(reads)

    def break_reads_into_kmers(self, reads):
        for read in reads:
            for j in range(len(read) - self.k + 1):
                yield read[j : j + self.k]

    def build_de_bruijn_graph(self, reads):
        for kmer in self.break_reads_into_kmers(reads):
            left, right = kmer[:-1], kmer[1:]
            if left == right:
                continue
            self.graph.setdefault(left, [[], 0])
            self.graph.setdefault(right, [[], 0])
            if right not in self.graph[left][0]:
                self.graph[left][0].append(right)
                self.graph[right][1] += 1

    def remove_edge(self, u, v):
        self.graph[u][0].remove(v)
        self.graph[v][1] -= 1
        self.edges_removed += 1

    def find_and_remove_outgoing(self, current, depth):
        if len(self.graph[current][0]) > 1 or self.graph[current][1] > 1:
            return False

        if depth == self.threshold:
            return False

        if len(self.graph[current][0]) == 0:
            return True

        nxt = self.graph[current][0][0]
        if self.find_and_remove_outgoing(nxt, depth + 1):
            self.remove_edge(current, nxt)
            return True

        return False

    def find_and_remove_incoming(self, current, depth):
        if depth == self.threshold:
            return False

        if len(self.graph[current][0]) == 0 or self.graph[current][1] > 1:
            return True

        nxt = self.graph[current][0][0]
        if self.find_and_remove_incoming(nxt, depth + 1):
            self.remove_edge(current, nxt)
            return True

        return False

    def remove_tips(self):
        for u, value in self.graph.items():
            incoming = False
            if len(value[0]) == 1 and value[1] == 0:
                incoming = True
            elif len(value[0]) <= 1:
                continue

            condition = True
            while condition:
                condition = False
                for w in list(value[0]):
                    if incoming:
                        ok = self.find_and_remove_incoming(w, 0)
                    else:
                        ok = self.find_and_remove_outgoing(w, 0)
                    if ok:
                        self.remove_edge(u, w)
                        condition = True
                        break

        return self.edges_removed


def main():
    reads = []
    for line in sys.stdin:
        read = line.strip()
        if read:
            reads.append(read)
    k = 15
    if reads and len(reads[0]) - 1 < k:
        k = len(reads[0]) - 1
    print(TipRemoval(reads, k).remove_tips())


if __name__ == "__main__":
    main()
