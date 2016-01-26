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

Test
----

    $ make

    $ ./kbartcheck fixtures/kbart.txt
    {"incomplete embargo":1,"records":72057}

    $ ./kbartcheck -skip fixtures/kbart.txt
    {"records":72056}
