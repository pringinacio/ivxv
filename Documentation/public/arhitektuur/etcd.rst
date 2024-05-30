..  IVXV arhitektuur

Lisa - ETCD Andmemudel
======================

ETCD on võti-väärtus andmebaas, kus talletatakse e-hääletamise
sisendnimekirju, e-hääli ning statistikat.

Ringkondade nimekiri
--------------------

.. table:: Ringkondade nimekiri
   :widths: 30 35 35

   +----------------------------+---------------------------+-----------------------+
   | **Võti**                   | **Väärtus**               | **Näide**             |
   +============================+===========================+=======================+
   | /districts                 | Juurvõti                  |                       |
   +----------------------------+---------------------------+-----------------------+
   | /districts/<EHAK-district> | Seab valija               | /districts/05241      |
   |                            | EHAK-ringkond paarile     |                       |
   |                            | vastavusse                | 0000.1                |
   |                            | ringkonnaidentifikaatori, |                       |
   |                            | mis viitab /choices       |                       |
   |                            | harusse                   |                       |
   +----------------------------+---------------------------+-----------------------+
   | /districts/counties        | Väli counties             | {                     |
   |                            | ringkondade               |                       |
   |                            | nimekirjast               | "0068": [             |
   |                            |                           |                       |
   |                            |                           | "0809",               |
   |                            |                           |                       |
   |                            |                           | "0624"                |
   |                            |                           |                       |
   |                            |                           | ],                    |
   |                            |                           |                       |
   |                            |                           | "0784": [             |
   |                            |                           |                       |
   |                            |                           | "0524"                |
   |                            |                           |                       |
   |                            |                           | ],                    |
   |                            |                           |                       |
   |                            |                           | "0079": [             |
   |                            |                           |                       |
   |                            |                           | "0796"                |
   |                            |                           |                       |
   |                            |                           | ],                    |
   |                            |                           |                       |
   |                            |                           | "0793": [             |
   |                            |                           |                       |
   |                            |                           | "0793"                |
   |                            |                           |                       |
   |                            |                           | ]                     |
   |                            |                           |                       |
   |                            |                           | }                     |
   +----------------------------+---------------------------+-----------------------+
   | /districts/version         | Nimekirja                 | ["NIMESTE,NIMI,123456 |
   |                            | allkirjastajad            | 78912                 |
   |                            |                           | 2019-02-22T13:58:48Z" |
   |                            |                           | ]                     |
   +----------------------------+---------------------------+-----------------------+

Valikute nimekiri
-----------------

.. table:: Valikute nimekiri
   :widths: 30 35 35

   +------------------------+-----------------------+-----------------------+
   | **Võti**               | **Väärtus**           | **Näide**             |
   +========================+=======================+=======================+
   | /choices               | Juurvõti              |                       |
   +------------------------+-----------------------+-----------------------+
   | /choices/<district-id> | Ringkonnale           | /choices/0000.1       |
   |                        | district-id vastav    |                       |
   |                        | valikute nimekiri.    | {                     |
   |                        |                       |                       |
   |                        |                       | "Erakond 1":{         |
   |                        |                       |                       |
   |                        |                       | "0000.101":"Nimi      |
   |                        |                       | Nimeste1",            |
   |                        |                       |                       |
   |                        |                       | "0000.102":"Nimi      |
   |                        |                       | Nimeste2",            |
   |                        |                       |                       |
   |                        |                       | "0000.103":"Nimi      |
   |                        |                       | Nimeste3"             |
   |                        |                       |                       |
   |                        |                       | },                    |
   |                        |                       |                       |
   |                        |                       | "Erakond 2":{         |
   |                        |                       |                       |
   |                        |                       | "0000.104":"Nimi      |
   |                        |                       | Nimeste4",            |
   |                        |                       |                       |
   |                        |                       | "0000.105":"Nimi      |
   |                        |                       | Nimeste5"             |
   |                        |                       |                       |
   |                        |                       | },                    |
   |                        |                       |                       |
   |                        |                       | "Üksikkandidaadid":{  |
   |                        |                       |                       |
   |                        |                       | "0000.106":"Nimi      |
   |                        |                       | Nimeste6"             |
   |                        |                       |                       |
   |                        |                       | }                     |
   |                        |                       |                       |
   |                        |                       | }                     |
   +------------------------+-----------------------+-----------------------+
   | /choices/version       | Nimekirja             | ["NIMESTE,NIMI,123456 |
   |                        | allkirjastajad        | 78912                 |
   |                        |                       | 2019-02-22T13:58:59Z" |
   |                        |                       | ]                     |
   +------------------------+-----------------------+-----------------------+

Valijate nimekiri
-----------------

.. table:: Valijate nimekiri
   :widths: 30 35 35

   +---------------------------------+-----------------------+-----------------------+
   | **Võti**                        | **Väärtus**           | **Näide**             |
   +=================================+=======================+=======================+
   | /voters                         | Juurvõti              |                       |
   +---------------------------------+-----------------------+-----------------------+
   | /voters/<version-id>            | Valijanimekirja       | /voters/1             |
   |                                 | versiooni juurvõti    |                       |
   +---------------------------------+-----------------------+-----------------------+
   | /voters/<version-id>/<voter-id> | Valija                | /voters/1/12345678912 |
   |                                 | <EHAK-district> antud |                       |
   |                                 | nimekirjas            | 05241                 |
   +---------------------------------+-----------------------+-----------------------+
   | /voters/<version-id>/version    | Valijatenimekirja     | ["NIMESTE,NIMI,123456 |
   |                                 | versioon              | 78912                 |
   |                                 |                       | 2019-02-22T13:58:59Z" |
   |                                 |                       | ]                     |
   +---------------------------------+-----------------------+-----------------------+
   | /voters/version                 | Aktuaalse             | 1                     |
   |                                 | valijatenimekirja     |                       |
   |                                 | versiooni ID          |                       |
   +---------------------------------+-----------------------+-----------------------+
   | /voters/previous                | Eelmise               | 0                     |
   |                                 | valijanimekirja       |                       |
   |                                 | versiooni ID          |                       |
   +---------------------------------+-----------------------+-----------------------+

Talletatud e-hääl
-----------------

.. table:: Talletatud e-hääl
   :widths: 30 35 35

   +-------------------------+------------------------+-----------------------+
   | **Võti**                | **Väärtus**            | **Näide**             |
   +=========================+========================+=======================+
   | /vote                   | Juurvõti               |                       |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>         | Hääle juurvõti,        | 16 baiti binaarandmed |
   |                         | unikaalne              |                       |
   |                         | identifikaator         |                       |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>/count   | Hääle kontrollimise    | 0                     |
   |                         | loendur                |                       |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>/ocsp    | Kehtivuskinnitus       | DER kodeeringus OCSP  |
   |                         |                        | vastus                |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>/time    | Hääle talletamise      | 2019-03-03T10:56:28.8 |
   |                         | kellaaeg               | 99925926Z             |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>/tspreg  | Registreerimiskinnitus | DER kodeeringus PKIX  |
   |                         |                        | ajatempel             |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>/type    | Allkirjastatud hääle   | BDOC                  |
   |                         | konteineri tüüp        |                       |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>/version | Hääle andmisel         | 0                     |
   |                         | kehtinud valijate      |                       |
   |                         | nimekirja versiooni    |                       |
   |                         | ID                     |                       |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>/vote    | E-hääl allkirjastatud  | BDOC vormingus        |
   |                         | konteineris            | allkirjastatud hääl   |
   +-------------------------+------------------------+-----------------------+
   | /vote/<vote-id>/voter   | Valija isikukood       | 12345678912           |
   +-------------------------+------------------------+-----------------------+

Statistikaliidesed
------------------

.. table:: Statistikaliidesed
   :widths: 30 35 35

   +-------------------------------+-----------------------+-----------------------+
   | **Võti**                      | **Väärtus**           | **Näide**             |
   +===============================+=======================+=======================+
   | /votes                        | Juurvõti              |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /votes/order                  | Hääletamisfaktide     |                       |
   |                               | järjestuse juurvõti   |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /votes/order/<seq>            | Konkreetse            | /votes/order/1        |
   |                               | hääletamisfakti       |                       |
   |                               | juurvõti              |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /votes/order/<seq>/admincode  | Hääletamisfaktiga     | 0796                  |
   |                               | seotud EHAK           |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /votes/order/<seq>/district   | Hääletamisfaktiga     | 10                    |
   |                               | seotud ringkonna      |                       |
   |                               | number                |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /votes/order/<seq>/voterid    | Hääletaja isikukood   | 12345678901           |
   |                               |                       |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /votes/order/<seq>/votername  | Hääletaja nimi        | NIMI NIMESTE          |
   |                               |                       |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /votes/stats                  | Viimase               | 12                    |
   |                               | hääletamisfakti       |                       |
   |                               | järjekorranumber      |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /voted                        | Juurvõti              |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /voted/latest                 | Viimati antud häälte  |                       |
   |                               | indeksi juurvõti      |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /voted/latest/<voter-id>      | Hääletaja poolt       | /voted/latest/1234567 |
   |                               | viimati antud hääle   | 8901                  |
   |                               | aeg ja identifikaator |                       |
   |                               | binaarkujul           | <2019-03-03T12:15:59Z |
   |                               |                       | ><vote-id>            |
   +-------------------------------+-----------------------+-----------------------+
   | /voted/stats                  | Jaoskonnapõhise       |                       |
   |                               | statistika indeksi    |                       |
   |                               | juurvõti              |                       |
   +-------------------------------+-----------------------+-----------------------+
   | /voted/stats/<voter-id>       | Hääle andmise         | /voted/stats/12345678 |
   |                               | kellaaeg koos         | 901                   |
   |                               | jaoskonnainfoga       |                       |
   |                               |                       | <0796><2019-02-22T14: |
   |                               |                       | 17:23Z>               |
   +-------------------------------+-----------------------+-----------------------+
   | /votes/voter/stats/<voter-id> | Tühi baitide massiiv  | /votes/voter/stats/   |
   |                               | (kasutatakse väärtuse | 394091044211          |
   |                               | versiooni, ning mitte |                       |
   |                               | väärtust ennast)      |                       |
   +-------------------------------+-----------------------+-----------------------+

Hääletamisseansid
-----------------

.. table:: Hääletamisseansid
   :widths: 30 35 35

   +-----------------------+-----------------------+-----------------------+
   | **Võti**              | **Väärtus**           | **Näide**             |
   +=======================+=======================+=======================+
   | /session              | Juurvõti              |                       |
   +-----------------------+-----------------------+-----------------------+
   | /session/<session-id> | RPC meetod, mis       | /session/0149468d2866 |
   |                       | kutsus antud          | 6fced7d73b32cc16225d  |
   |                       | funktsiooni välja +   |                       |
   |                       | x1F + kasutaja        |                       |
   |                       | autentimismeetod      |                       |
   +-----------------------+-----------------------+-----------------------+
