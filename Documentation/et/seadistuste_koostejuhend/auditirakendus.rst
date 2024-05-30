..  IVXV kogumisteenuse haldusteenuse kirjeldus

.. _app-auditor:

Auditirakendus
==============

IVXV võtmerakendus võimaldab kasutada tõestatavat dekrüpteerimist -
koos tulemusega väljastatakse lugemistõend e-häälte korrektse avamise
kohta. Vältimaks häälte salajasuse rikkumist lugemistõendi kontrollil
võimaldab IVXV kasutada häälte miksimist, mis säilitab häälte sisu
kuid eemaldab krüptograafiliselt seose konkreetse hääle ja selle hääle
andnud isiku vahel.

IVXV kasutab e-häälte miksimiseks tarkvara Verificatum, mis võtab
sisendiks krüpteeritud hääled ning annab väljundiks miksitud
krüpteeritud hääled ja miksimistõendi.

Miksimistõendi ja lugemistõendi kontroll toimub auditirakenduse tööriistadega
*convert*, *mixer* ja *decrypt*.

Töötlemisrakenduse toimingute kontroll toimub auditirakenduse tööriistaga
*integrity*.

#. Tööriist *convert* kontrollib, teisenduste korrektsust IVXV
   andmevormingute ja Verificatumi andmevormingute vahel.
#. Tööriist *mixer* kontrollib miksimistõendi korrektsust.
#. Tööriist *decrypt* kontrollib lugemistõendi korrektsust.
#. Tööriist *integrity* kontrollib töötlemisrakenduse logide korrektsust.

E-häälte korrektse kokkulugemise kontrolliks on vajalik ja piisav
kasutada kõiki nelja auditirakenduse tööriista.

Kõigi tööriistade kasutamine eeldab allkirjastatud usaldusjuure ja
konkreetse tööriista seadistuste olemasolu. Alljärgnevalt kirjeldame
konkreetsete tööriistade seadistusi.

.. _auditor-convert:

E-häälte korrektse teisendamise kontroll
----------------------------------------

Verificatumi poolt koostatud miksimistõendi formaat on erinev IVXV raamistikus
kasutatavast formaadist, samuti erinevad IVXV ning Verificatumi
krüpteeritud häälte formaadid. IVXV raamistikku on pakendatud
adapterid formaaditeisendusteks, auditirakendus pakub võimalust nende
teisenduste korrektsuse kontrolliks.

Tööriist *convert* kontrollib, et Verificatumi poolt väljastatud
miksimistõend vastab failidele IVXV raamistikus.

:convert.input_bb: IVXV miksimiseelse e-valimiskasti asukoht.

:convert.output_bb: IVXV miksimisjärgse e-valimiskasti asukoht.

:convert.pub: IVXV avaliku võtme asukoht.

:convert.protinfo: Verificatumi miksimise protokollifaili asukoht.

:convert.proofdir: Verificatumi miksimistõendi asukoht.

:file: `auditor.convert.yaml`:

.. literalinclude:: config-examples/auditor.convert.yaml
   :language: yaml
   :linenos:

.. _auditor-mix:

E-häälte miksimistõendi kontroll
--------------------------------

Tööriist *mixer* kontrollib Verificatumi miksimistõendi korrektsust.

:mixer.protinfo: Verificatumi miksimistõendi protokollifaili asukoht.

:mixer.proofdir: Verificatumi miksimistõendi asukoht.

:mixer.threaded: Kasuta mitmelõimelist implementatsiooni. Vaikimisi
                 väärtus on väär. Kasutatavate lõimede arv sõltub
                 käsurea-argumentidest. Käsurea-argumentide puudumise
                 korral valitakse optimaalne lõimede arv lähtudes
                 tuvastatud tuumade arvust.

:file:`auditor.mixer.yaml`:

.. literalinclude:: config-examples/auditor.mixer.yaml
   :language: yaml
   :linenos:

.. _auditor-decrypt:

E-häälte lugemistõendi kontroll
-------------------------------

Tööriist *decrypt* kontrollib lugemistõendi korrektsust.

:decrypt.proofs: Kehtivate sedelite lugemistõendi asukoht.

:decrypt.pub: Dekrüpteerimiseks kasutatud salajasele võtmele vastava avaliku
              võtme asukoht.

:decrypt.discarded: Loend kehtetutest sedelitest.

:decrypt.anon_bb: Töötlemisrakenduse või miksimisrakenduse poolt loodud
                  e-valimiskast anonüümistatud häältega.

:decrypt.plain_bb: Dekrüpteeritud valimiskast.

:decrypt.tally: Elektroonilise hääletamise tulemus.

:decrypt.candidates: Valimise valikute nimekiri allkirjastatud kujul.

:decrypt.districts: Valimise ringkondade nimekiri allkirjastatud kujul.

:decrypt.out: Lugemistõendi kontrolli tulemuste asukoht. Tegemist on
              kataloogiga kuhu salvestatakse sedelid, mille
              lugemistõend oli kehtetu.

:decrypt.invalidity_proofs: Valikuline kehtetute sedelite lugemistõendi
                            asukoht.

:decrypt.abort_early: Valikuline auditirakenduse peatamine esimese läbikukutud
                      kontrolli korral. Vaikimisi väärtus on tõene.

:file:`auditor.decrypt.yaml`:

.. literalinclude:: config-examples/auditor.decrypt.yaml
   :language: yaml
   :linenos:

.. _auditor-integrity:

Töötlemisrakenduse logide kontroll
----------------------------------

Tööriist *integrity* kontrollib, et töötlemisrakenduse poolt väljastatud
logid ühendavad e-valimiskasti anonüümitud e-valimiskastiga.

:integrity.ballotbox: Kogumisteenusest väljastatud e-valimiskast.

:integrity.anon_bb: Töötlemisrakenduse poolt loodud e-valimiskast
                    anonüümistatud häältega.

:integrity.log_accepted: Vastuvõetud häälte *log1* fail.

:integrity.log_squashed: Tühistatud korduvhäälte häälte *log2* fail.

:integrity.log_revoked: Jaoskonnainfo põhjal tühistatud ja ennistatud häälte
                        *log2* fail.

:integrity.log_anonymised: Lugemisele läinud häälte *log3* fail.

:integrity.bb_errors: E-valimiskasti töötlemisvigade raport.

:integrity.abort_early: Valikuline auditirakenduse peatamine esimese
                        läbikukutud kontrolli korral. Vaikimisi
                        väärtus on tõene.

:file: `auditor.integrity.yaml`:

.. literalinclude:: config-examples/auditor.integrity.yaml
   :language: yaml
   :linenos:
