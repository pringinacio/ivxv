..  IVXV protokollid

=========================
E-valimiskasti töötlemine
=========================


Tühistus- ja ennistusnimekiri
-----------------------------

Tühistus- ja ennistusnimekiri sisaldab andmeid isikute kohta, kelle e-hääl tuleb tühistada (ei lähe arvesse valimistulemuste kokkulugemisel) või ennistada (s.t. tühistatakse eelnev tühistamine ning häälte uuesti üle lugemisel võetakse ennistatud e-hääl arvesse). Nimekiri laaditakse süsteemi digitaalselt allkirjastatud dokumendina, mille andmefaili vorming on järgmine:

.. literalinclude:: ../../common/schema/ivxv.revoke.schema
   :language: json
   :linenos:

Näide:

.. literalinclude:: ../../common/schema/ivxv.revoke.schema.example
   :language: json


E-hääletanute nimekiri
----------------------

E-hääletanute nimekiri on pärast e-hääletamise lõppu väljastatav nimekiri
e-hääletanud isikutest, sordituna valimisjaoskondade kaupa. Dokument
genereeritakse töötlemisrakenduse poolt.

.. literalinclude:: ../../common/schema/ivxv.voterlist.schema
   :language: json
   :linenos:


Näide:

.. literalinclude:: ../../common/schema/ivxv.voterlist.schema.example
   :language: json


Hääletamistulemus
-----------------

Võtmerakenduse poolt dekrüpteeritud ning summeeritud hääled jagatud
valimisringkondade ja jaoskondade kaupa.

Hääletamistulemuste failis peavad iga jaoskonna kohta olema järgmised andmed.

#. Rikutud ja kehtetute häälte arvu näitav kirje. Seda ka juhul, kui
   valimisjaoskonnas polnud ühtki rikutud või kehtetut häält: sellisel juhul on
   häälte arv null.

#. Iga valiku poolt antud häälte arvu näitav kirje. Seda ka juhul, kui
   valimisjaoskonnas ei antud selle valiku poolt ühtki häält: sellisel juhul on
   häälte arv null.

.. literalinclude:: ../../common/schema/ivxv.result.schema
   :language: json
   :linenos:

Näide:

.. literalinclude:: ../../common/schema/ivxv.result.schema.example
   :language: json



E-valimiskast
-------------

Fail sisaldab kogumisteenuse poolt vastu võetud hääli koos häälte juurde
kuuluvate andmetega.

Faili vorming on Zip64 konteiner.

Valija-spetsiifilised kaustad asuvad vahetult juurkausta `votes` all.

Faili sisu:

* :file:`votes/<voter id>/`

 * :file:`<timestamp>.version`

 * :file:`<timestamp>.<vote type>`

 * :file:`<timestamp>.<qualifier>*`

kus:

* ``<voter id>`` on valija identifikaator, Eesti puhul isikukood;

* ``<timestamp>`` on hääle esitamise kellaaeg vormingus
  ``yyyymmddhhmmssmmm±zzzz``;

  * see kellaaeg kajastab hetke, mil päring kogumisteenusesse tehti, ja on
    antud lihtsalt valimiskasti inimloetavuse parandamiseks; hääle tegelik ajamärk või
    -tempel on mõne kvalifitseeriva vastuse sees;

* ``<vote type>`` on valikute konteineri tüüp, Eesti puhul BDOC;

  * kusjuures BDOC ise on lihtsalt põhiprofiiliga ja ei sisalda
    kvalifitseerivad parameetreid (kehtivuskinnitusi, ajamärgendeid,
    ajatempleid),

* ``<qualifier>`` on häält kvalifitseeriva protokolli tüüp, millest hetkel
  võimalikud on:

   * ``ocsp`` - *Online Certificate Status Protocol* (kehtivuskinnitus, `RFC
     6960 <https://tools.ietf.org/html/rfc6960>`_) kinnitab valija
     allkirjastamissertifikaadi kehtivust hääle andmise hetkel,

   * ``tsp`` - *Time-Stamp Protocol* (ajatempel, `RFC 3161
     <https://tools.ietf.org/html/rfc3161>`_) kinnitab, et päringu tegemise
     hetkeks oli hääl olemas,

   * ``tspreg`` - sama, mis ``tsp``, aga nonsiks pannakse kogumisteenuse
     allkiri päringu ``MessageImprint`` elemendil, et häält registreerida.

* Iga hääle kohta esinevad failid on:

 * ``<timestamp>.version`` - hääle andmise ajal kehtinud valijate nimekirja
   versioon;

 * ``<timestamp>.<vote type>`` - valikute konteiner, mille sees on valiku
   identifikaator kujul ``<valimise id>.<küsimuse id>.ballot``. Eesti puhul
   BDOC-konteineris olev vastava nimega fail;

 * ``<timestamp>.<qualifier>`` - häält kvalifitseeriva protokolli päringu
   vastus; neid võib esineda mitu, aga iga protokolli kohta maksimaalselt
   üks.


Anonüümistatud e-valimiskast
----------------------------

Valimisringkondade ja jaoskondade järgi grupeeritud krüpteeritud hääled.
Anonüümistatud e-valimiskastis puudub informatsioon valijate kohta.

Anonüümistatud e-valimiskast on töötlemisrakenduse väljund ning võtmerakenduse
dekrüpteerimise tööriista sisend.

.. literalinclude:: ../../common/schema/ivxv.anon-bb.schema
   :language: json
   :linenos:

Näide:

.. literalinclude:: ../../common/schema/ivxv.anon-bb.schema.example
   :language: json
