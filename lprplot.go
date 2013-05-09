// lprplot.go

/*
This package mimics an old-school Line Printer Plot

  Install
  -------
  go get github.com/hotei/lprplot

  Features
  --------
  Pretty good imitation of LPR output in about 100 loc.
  See dynamo for example of usage and output
  
  Limitations
  -----------
  Overlaps somewhat with Features :-)
  Requires caller to set min/max ranges by symbol which may require a pass over data 
    This could be avoided, but at some cost in complexity of plot
    
  TODO
  ----
  * TBD
  
 (c) 2013 David Rook - License is BSD style - see LICENSE.md
  Also see README.md for more info
  
*/
package lprplot

import(
"fmt"
"os"
)

// holds lower/upper limits for plotting lines
type minmax struct {
	minimum float64
	maximum float64
}

type Lplot struct {
	width      int			// page width, 80 or 132 were common
	bufr       []byte		// compose buffer
	edge       int			// column where axis occurs
	out        * os.File 	//  io.Writer better?
	plot_range map[byte]*minmax 	// each symbol has lo-hi range associated
}

func (p *Lplot) Init(w int, f *os.File) {
	p.width = w
	p.bufr = make([]byte, p.width)
	for i := 0; i < p.width; i++ {
		p.bufr[i] = ' '
	}
	p.edge = 9
	p.plot_range = make(map[byte]*minmax, 10)
	p.out = f
}

// presets can clip if narrower than actual data
func (p *Lplot) SetRange(symb byte, mini, maxi float64) {
	//	fmt.Printf("set plot_range = %v\n",p.plot_range)
	var m = minmax{mini, maxi}
	p.plot_range[symb] = &m
}

// prepare plot for next print line by printing next year and left side of box
func (p *Lplot) AxisLabel(s string) {
	ss := fmt.Sprintf("%7s | ", s)
		copy(p.bufr[0:p.edge+1],ss)
	return
	
	for ndx, b := range s {
		p.bufr[ndx] = byte(b)
	}
}

// compose the line by adding one symbol, doesn't create output (Advance does)
func (p *Lplot) Plot(sym byte, val float64) {
	rmax := p.plot_range[sym].maximum
	rmin := p.plot_range[sym].minimum
	step := (rmax - rmin) / float64(p.width-p.edge)
	i := int((val - rmin) / step)
	if i < 0 {
		i = 0
	}
	pt := i + p.edge
	if pt >= p.width {
		// fmt.Printf("cant happen - value %g translates to pt[%d]", val, pt)
		return
	}
	if p.bufr[pt] == ' ' {
		p.bufr[pt] = sym
	} else {
		p.bufr[pt] = '*'	// multiple symbols at same location
	}
}

// output the line to whatever Writer is given, also to screen
func (p *Lplot) Advance() {
	//	fmt.Printf("len_bufr(%d)\n",len(p.bufr))
	for i := 0; i < p.width; i++ {
		fmt.Printf("%c", p.bufr[i])
		n, err := p.out.Write([]byte{p.bufr[i]})
		if err != nil || n != 1 {
			fmt.Printf("write to plotfile failed\n")
			os.Exit(-2)
		}
		p.bufr[i] = ' '
	}
	fmt.Printf("|\n")
	n, err := p.out.Write([]byte{'\n'})
	if err != nil || n != 1 {
		fmt.Printf("write to plotfile failed\n")
		os.Exit(-2)
	}
}

// print non-plot info like legends, extra lines for pagination etc.
func (p *Lplot) WriteString(s string) {
	p.out.WriteString(s)	
}