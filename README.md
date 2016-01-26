README
======

Holdings files are used to specify the availability of a resources. This
package provides support for various holding file formats and a common
interface.

Supported formats:

* Google
* KBART
* OVID

Not supported:

* http://www.loc.gov/marc/holdings/echdhome.html

Testdrive
---------

    $ make

    $ kbartcheck fixtures/kbart.txt
    {"incomplete embargo":1,"records":72057}

    $ kbartcheck -skip fixtures/kbart.txt
    {"records":72056}

    $ holdingcov -issn 1325-9210 -file fixtures/kbart.txt -date 2009-10-10
    0   OK  No restrictions.
    1   NO  Not covered: after coverage interval

    $ holdingcov -issn 1520-4898 -date 1995 -volume 29 -file fixtures/kbart.txt
    0   NO  Not covered: before coverage interval
    1   NO  Not covered: after coverage interval

    $ holdingcov -issn 1520-4898 -date 1995 -volume 28 -file fixtures/kbart.txt
    0   NO  Not covered: before coverage interval
    1   OK  No restrictions.

    $ holdingcov -issn 1613-4141 -date 2015 -volume 1 -issue 2 -file fixtures/kbart.txt
    0   OK  No restrictions.
    1   NO  Moving wall applies.
    2   NO  Not covered: after coverage interval

    $ holdingcov -issn 1613-4141 -date 2015 -volume 1 -issue 2 -file fixtures/ovid.xml -format ovid
    0   OK  No restrictions.
    1   OK  No restrictions.
    2   NO  Not covered: after coverage interval

    $ holdingcov -issn 1613-4141 -date 2015 -volume 1 -issue 2 -file fixtures/google.xml -format google
    0   OK  No restrictions.
    1   OK  No restrictions.
    2   NO  Not covered: after coverage interval

    $ make clean
