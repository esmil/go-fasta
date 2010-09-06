/*
 * This file is part of go-fasta.
 *
 * go-fasta is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of
 * the License, or(at your option) any later version.
 *
 * go-fasta is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with go-fasta.  If not, see <http://www.gnu.org/licenses/>.
 */
package fasta

import (
	"testing"
	"fmt"
	"os"
)

var str string = ";hej\n\r>> dav\rAb,C:-djJjJEf\ngH"

func TestParser1(t *testing.T) {
	var p Parser
	var f FASTA

	if done, rest := p.Feed([]byte(str)); done || rest != nil {
		t.Errorf("Weird return: %v, %v", done, rest)
		return
	}

	p.Result(&f)
	fmt.Printf("%#v\n%s", &f, &f)
}

func TestParseOne(t *testing.T) {
	var f FASTA

	filename := "laminin.fasta"
	file, err := os.Open(filename, os.O_RDONLY, 0)
	if err != nil {
		t.Errorf("Error opening %s: %s", filename, err)
		return
	}
	defer file.Close()

	err = ParseOne(file, &f)
	if err != nil {
		t.Errorf("Error parsing %s: %s", filename, err)
		return
	}

	fmt.Print(&f)
}

func TestParseAll(t *testing.T) {
	filename := "wikiall.fasta"
	file, err := os.Open(filename, os.O_RDONLY, 0)
	if err != nil {
		t.Errorf("Error opening %s: %s", filename, err)
		return
	}
	defer file.Close()

	fasta, err := ParseAll(file)
	if err != nil {
		t.Errorf("Error parsing %s: %s", filename, err)
		return
	}

	for i := range fasta {
		fmt.Print(fasta[i])
	}
}
