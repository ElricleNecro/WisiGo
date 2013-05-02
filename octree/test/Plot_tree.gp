#!/usr/local/bin/gnuplot -persist
#
#    
#    	G N U P L O T
#    	Version 4.5 patchlevel 0
#    	last modified 2011-11-18 
#    	System: Linux 3.1.0-1-686-pae
#    
#    	Copyright (C) 1986-1993, 1998, 2004, 2007-2011
#    	Thomas Williams, Colin Kelley and many others
#    
#    	gnuplot home:     http://www.gnuplot.info
#    	mailing list:     gnuplot-beta@lists.sourceforge.net
#    	faq, bugs, etc:   type "help seeking-assistance"
#    	immediate help:   type "help"
#    	plot window:      hit 'h'
set terminal x11  noraise persist enhanced
# set output
unset clip points
set clip one
unset clip two
set bar 1.000000 front
unset border
set timefmt z "%d/%m/%y,%H:%M"
set zdata 
set timefmt y "%d/%m/%y,%H:%M"
set ydata 
set timefmt x "%d/%m/%y,%H:%M"
set xdata 
set timefmt cb "%d/%m/%y,%H:%M"
set timefmt y2 "%d/%m/%y,%H:%M"
set y2data 
set timefmt x2 "%d/%m/%y,%H:%M"
set x2data 
set boxwidth
set style fill  empty border
set style rectangle back fc lt -3 fillstyle   solid 1.00 border lt -1
set style circle radius graph 0.02, first 0, 0 
set style ellipse size graph 0.05, 0.03, first 0 angle 0 units xy
set dummy x,y
set format x "%g"
set format y "%g"
set format x2 "% g"
set format y2 "% g"
set format z "% g"
set format cb "% g"
set format r "% g"
set angles radians
unset grid
set raxis
set key title ""
set key bmargin left horizontal Right noreverse enhanced autotitles box linetype -1 linewidth 1.000
set key noinvert samplen 4 spacing 1 width 0 height 0 
set key maxcolumns 0 maxrows 0
set key noopaque
unset key
unset label
unset arrow
set style increment default
unset style line
unset style arrow
set style histogram clustered gap 2 title  offset character 0, 0, 0
unset logscale
set offsets 0, 0, 0, 0
set pointsize 1
set pointintervalbox 1
set encoding utf8
unset polar
unset parametric
unset decimalsign
set view 60, 30, 1, 1
set samples 100, 100
set isosamples 10, 10
set surface
unset contour
set clabel '%8.3g'
set macros
set mapping cartesian
set datafile separator whitespace
unset hidden3d
set cntrparam order 4
set cntrparam linear
set cntrparam levels auto 5
set cntrparam points 5
set size ratio 0 1,1
set origin 0,0
set style data points
set style function lines
set xzeroaxis linetype -2 linewidth 1.000
set yzeroaxis linetype -2 linewidth 1.000
set zzeroaxis linetype -2 linewidth 1.000
set x2zeroaxis linetype -2 linewidth 1.000
set y2zeroaxis linetype -2 linewidth 1.000
set ticslevel 0.5
set mxtics default
set mytics default
set mztics default
set mx2tics default
set my2tics default
set mcbtics default
set noxtics
set noytics
set ztics border in scale 1,0.5 nomirror norotate  offset character 0, 0, 0 autojustify
set ztics autofreq  norangelimit
set nox2tics
set noy2tics
set cbtics border in scale 1,0.5 mirror norotate  offset character 0, 0, 0 autojustify
set cbtics autofreq  norangelimit
set rtics axis in scale 1,0.5 nomirror norotate  offset character 0, 0, 0 autojustify
set rtics autofreq  norangelimit
set title "Fonctionnement du Tree Code\nNbMin = 5" 
set title  offset character 0, 0, 0 font "" norotate
set timestamp bottom 
set timestamp "" 
set timestamp  offset character 0, 0, 0 font "" norotate
set rrange [ * : * ] noreverse nowriteback
set trange [ * : * ] noreverse nowriteback
set urange [ * : * ] noreverse nowriteback
set vrange [ * : * ] noreverse nowriteback
set xlabel "" 
set xlabel  offset character 0, 0, 0 font "" textcolor lt -1 norotate
set x2label "" 
set x2label  offset character 0, 0, 0 font "" textcolor lt -1 norotate
set xrange [ * : * ] noreverse nowriteback
set x2range [ * : * ] noreverse nowriteback
set ylabel "" 
set ylabel  offset character 0, 0, 0 font "" textcolor lt -1 rotate by -270
set y2label "" 
set y2label  offset character 0, 0, 0 font "" textcolor lt -1 rotate by -270
set yrange [ * : * ] noreverse nowriteback
set y2range [ * : * ] noreverse nowriteback
set zlabel "" 
set zlabel  offset character 0, 0, 0 font "" textcolor lt -1 norotate
set zrange [ * : * ] noreverse nowriteback
set cblabel "" 
set cblabel  offset character 0, 0, 0 font "" textcolor lt -1 rotate by -270
set cbrange [ * : * ] noreverse nowriteback
set zero 1e-08
set lmargin  -1
set bmargin  -1
set rmargin  -1
set tmargin  -1
set locale "fr_FR.UTF-8"
set pm3d explicit at s
set pm3d scansautomatic
set pm3d interpolate 1,1 flush begin noftriangles nohidden3d corners2color mean
set palette positive nops_allcF maxcolors 0 gamma 1.5 color model RGB 
set palette defined ( 0 0 0 1, 0.5 0 1 0, 1 1 0 0 )
set colorbox default
set colorbox vertical origin screen 0.9, 0.2, 0 size screen 0.05, 0.6, 0 front bdefault
set style boxplot candles range  1.50 outliers pt 7 separation 1 labels auto unsorted
set loadpath 
set fontpath 
set psdir
set fit errorvariables
pdf( file ) = sprintf("set t push;set t pdf;set o '%s';replot;set o;set t pop", file )
png( file ) = sprintf("set t push;set t png;set o '%s';replot;set o;set t pop", file )
latex( file ) = sprintf("set t push;set t tex;set o '%s';replot;set o;set t pop", file )
tikz( file ) = sprintf("set t push;set t tikz nostandalone color solid;set o '%s';replot;set o;set t pop", file )
stikz( file ) = sprintf("set t push;set t tikz standalone color solid;set o '%s';replot;set o;set t pop", file )
GNUTERM = "wxt"
GPFUN_pdf = "pdf( file ) = sprintf(\"set t push;set t pdf;set o '%s';replot;set o;set t pop\", file )"
GPFUN_png = "png( file ) = sprintf(\"set t push;set t png;set o '%s';replot;set o;set t pop\", file )"
GPFUN_latex = "latex( file ) = sprintf(\"set t push;set t tex;set o '%s';replot;set o;set t pop\", file )"
GPFUN_tikz = "tikz( file ) = sprintf(\"set t push;set t tikz nostandalone color solid;set o '%s';replot;set o;set t pop\", file )"
GPFUN_stikz = "stikz( file ) = sprintf(\"set t push;set t tikz standalone color solid;set o '%s';replot;set o;set t pop\", file )"
plot 'Particule-tc.dat' u 1:2:3 w l palette lt 0, 'Particule-cl.dat' u 1:2:3 w p palette pt 3 ps 2
plot 'Particule-tc.dat' u 1:2 w l lt 0, 'Particule-cl.dat' u 1:2:4 w p palette pt 3 ps 2 , 'Particule-af.dat' u 1:2 w p pt 8 ps 4
#    EOF
