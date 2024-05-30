# E-hääletamiste nimekiri

kavand v0.4

## Muutelugu

- Eemaldatud sõna "jooksev" / Priit Parmakson, 23.11.2022
- Lisatud OpenAPI kirjelduse publitseerimise teave. / Priit Parmakson, 13.12.2022
- Täpsustatud e-hääletamiste väljanäitamist VIS3-s. / Priit Parmakson, 15.11.2022
- Lisatud jõudluskaalutlused. / Priit Parmakson, 01.11.2022
- Asendatud "e-hääletanu" -> "e-hääletamine". Põhjus: inimene võib e-hääletada mitu korda (ümberhääletamine); vastavalt korrigeeritud JSON väljanimesid "evotersbatch" -> "evotingsbatch". /  Priit Parmakson, 01.09.2022
- Eemaldatud: hääletamise kuupäeva/aja edastamine. Põhjus: turvakaalutlus (häälemüügi takistamine). / Priit Parmakson, 01.09.2022 (Sven Heibergiga arutelu alusel)
- Lisatud: 1) ärivajaduse täpsustus; 2) e-hääletamise fakti aja edastamine; 3) edastuse kontroll ja vajadusel täiendamine või parandamine lõpliku faili abil. / Priit Parmakson, 04.05.2022
- Arutelu - RVT ja RIA inimesed, 04.05.2022
- "E-hääletanute jooksev nimekiri", kavand v0.2. / Priit Parmakson, 02.05.2022
- Arutelu: X-tee kasutamine, JSON - Tarmo Hanga, Priit Parmakson, apr 2022
- "EHS- VIS3 liidestus" (kavand) - Indrek Leesi, apr 2022

## Ülevaade

Käesolev spetsifikatsioon määratleb protokolli e-hääletamiste nimekirja edastamiseks e-hääletamise süsteemist (edaspidi - EHS) Valimiste infosüsteemi (edaspidi - VIS3).

Varasematel valimistel (k.a KOV 2021) on peale hääletamise lõppu, valimispäeval edastatud EHS-st VIS3-le e-hääletanute nimekiri (edaspidi - e-hääletanute lõplik nimekiri). See edastus on spetsifitseeritud: [E-hääletanute nimekiri](https://github.com/e-gov/VIS3-EHS/blob/main/4_e_haaletanute_nimekiri/SPEC.md).

Siiski on vajadus VIS3-s saada teavet, kas valija on e-hääletanud, juba enne ülalnimetatud lõplikku edastust. Valija, kes on e-hääletanud, võib tulla eelhääletamise perioodil valimisjaoskonda ja soovida paberil hääletada. Valimisjaoskonna töötajal oleks hea omada võimalust VIS3-st vaadata, kas valija on e-hääletanud. Seda teavet saab muuhulgas kasutada valija hoiatamiseks, et paberhääletamisega tema e-hääletamine tühistub.

E-hääletamiste nimekiri edastatakse EHS-st VIS3-e X-tee teenusega.

Teenust pakub EHS. VIS3 pöördub regulaarselt teenuse poole. EHS edastab e-hääletamiste andmete paki. VIS3 saab paki ja salvestab andmed VIS3 andmebaasi.

Sünkroonimiseks kasutatakse e-hääletamiste nummerdamist.

## Kaalutlused

Teenus on arendatud X-tee REST võimalusi kasutades, vastavalt X-tee REST sõnumiprotokollile ([X-Road: Message Protocol for REST](https://www.x-tee.ee/docs/live/xroad/pr-rest_x-road_message_protocol_for_rest.html)). Andmed väljastatakse JSON-vormingus.

X-teed on otstarbekas kasutada, sest kuigi VIS3 ja EHS võivad olla käitatud samas taristus, on ikkagi vaja tagada usaldus, logimine ja paindlikkus - omadused, mille tagamine X-tee kasutamisega kokkuvõttes tõenäoliselt ei saaks olema ei lihtsam ega odavam.

Protokollis on arvesse võetud Rahvastikuregistri ja VIS3 vaheliste X-tee teenuste kasutamise kogemust (REST sõnumiprotokoll, JSON, OpenAPI).

Erilist tähelepanu on pööratud andmete re-sünkroonimise võimalusele tõrgete korral. Selleks on pakkide pärimine kavandatud idenmpotentsena.

E-hääletanute nimekirja "peegeldamiseks" EHS-st VIS3-e võiks põhimõtteliselt kasutada ka mõnda standardset sünkroonimisprotokolli (nt Git, Rsync vms, vt \[1], \[2], \[3]). Kuna vajadus on suhteliselt lihtne, siis seda ei ole tehtud.

Jõudluskaalutlused: E-hääletamise jooksvat nimekirja edastatakse e-hääletamise perioodil (6 päeva). E-hääletamise fakte edastatakse kokku u 300 000. E-hääletamise perioodi esimesel tunnil võib oodata u 5000 e-hääletamist; perioodi viimasel tunnil u 10 000 e-hääletamist.

## E-hääletamiste järjenumbrid

Tagamaks, et e-hääletamiste nimekiri kantakse EHS-st VIS3-e õigeaegselt ja täielikult, kasutatakse järjenumbreid. EHS omistab igale e-hääletamisele järjenumbri (ingl Sequence Number).

Järjenumber on naturaalarv, alates ühest. Väärtus `0` tähistab olukorda, kus e-hääletamisi veel ei ole toimunud.

Igas valimissündmuses on oma numeratsioon.

## Teenus

Teenuse vastutav töötleja on Riigi Valimisteenistus (RVT). Teenuse volitatud töötleja on Riigi Infosüsteemi Amet (RIA). Teenust osutav süsteem on EHS.

Teenust kasutav süsteem on VIS3. Teenust kasutava süsteemi vastutav töötleja on RIA.

Teenuse ärinimi on "E-hääletamiste nimekiri". Teenuse tehniline nimi (X-tee REST teenusekood, Service Code) on `e-votings-running-list`.

Teenus pakub järgmisi otspunkte:

1  `GET /elections`.  "Valimissündmuste loetelu" väljastab aktiivsete valimissündmuste loetelu. Aktiivne valimissündmus teenuse kontekstis on selline, mille kohta EHS on valmis väljastama e-hääletamiste nimekirja.

2  `GET /elections/{electionId}/lastseqno`. "Viimane järjenumber" väljastab konkreetse valimissündmuse viimase EHS-s registreeritud e-hääletamise järjenumbri.

3  `GET /elections/{electionId}/evotingsbatchfrom/{fromseqno}`. "e-hääletamiste pakk". Selle päringuga pärib VIS3 EHS-lt valimissündmuse `{electionId}` e-hääletamiste paki, alatest e-hääletamisest järjenumbriga `{fromseqno}`.

Päringute ja vastuste andmestruktuuride ja samuti vastuskoodide spetsifikatsiooni vt OpenAPI spetsifikatsioonis: [ehs-xroad-api.yaml](ehs-xroad-api.yaml).

Teenuse OpenAPI spetsifikatsioon publitseeritakse turvaserveris, vastavalt spetsifikatsioonile [https://www.x-tee.ee/docs/live/xroad/pr-mrest_x-road_service_metadata_protocol_for_rest.html#4-retrieving-list-of-services](X-Road: Service Metadata Protocol for REST), jaotis 5 "Retrieving the OpenAPI description of a Service".

Teenuse OpenAPI spetsifikatsioon näitepäring:

`curl -H "accept: application/json" -H "X-Road-Client:ee-dev/COM/<MEMBER>/dev" "https://.../r1/ee-dev/COM/<MEMBER>/ehs/getOpenAPI?serviceCode=e-votings-running-list"`

## Otspunkt "Valimissündmuste loetelu"

Näide.

Päring: `GET /elections`

VIS3 pärib EHS-lt aktiivsete valimissündmuste loetelu.

Vastus:

```
{ "elections": [ { "name": "RK_2023" } ] }
```

EHS vastab, et aktiivseid valimissündmusi on üks - `RK_2023`.

Kui aktiivseid valimissündmusi ei ole, siis EHS peab vastuses saatma tühja massiivi (JSON Array).

## Otspunkt "Viimane järjenumber"

Näide.

Päring: `GET /elections/RK_2023/lastseqno`

VIS3 pärib valimissündmuse `RK_2023` viimase e-hääletamise järjenumbrit.

Vastus:

```
{
  "electionname": "RK_2023",
  "lastseqno": 54002
}
```

EHS vastab, et valimissündmuse `RK_2023` viimase e-hääletamise järjenumber on `54002`.

Kui valimissündmus on EHS-le tundmatu, siis EHS vastab HTTP vastuskoodiga `404 Not Found`.

Järjenumbrid algavad ühest (`1`). Kui valimisündmuses ei ole veel ükski valija e-hääletanud, siis vastab EHS `lastseqno` väärtusega `0`.

## Otspunkt "e-hääletanute pakk"

Näide 3.

Päring: `G /elections/RK_2023/evotingsbatchfrom/54001`

VIS3 pärib valimissündmuse `RK_2023` e-hääletamiste andmeid, alates järjenumbrist `54001`.

Vastus:

```
{
  "electionname": "RK_2023",
  "fromseqno": 54001,
  "batchmaxsize": 100,
  "evotingsbatch": [
    {
      "seqno": 54001,
      "idcode": "38101010020",
      "votername": "LEO KASS",
      "kovcode": "0068",
      "electoraldistrictno": 4
    },
    {
      "seqno": 54002,
      "idcode": "38101010021",
      "votername": "MARK KOER",
      "kovcode": "0305",
      "electoraldistrictno": 4
    }
  ]

}
```

EHS vastab, et saadab valimissündmuse `RK_2023` e-hääletamiste andmeid, alates järjenumbrist `54001`, pakina, milles on kuni `100` kirjet. Konkreetses pakis on kaks kirjet, kuna paki moodustamise hetkel on viimase EHS-s registreeritud e-hääletamise järjenumber `54002`. Esimene kirje tähendab, et valija `LEO KASS`, isikukoodiga `38101010020` on valimistel `RK_2023` e-hääletanud. Kirjes on ka valija KOV EHAK-kood ja valimisringkonna number.

E-hääletamise fakti kohta EHS e-hääletamise aega ei saada - turvakaalutlustel.

Paki maksimaalsuuruse `batchmaxsize` määrab EHS, arvestusega, et andmed saadetakse X-tee vastussõnumi kehas (mitte manuses). Vastussõnumi töötlemisel turvaserveris loetakse keha üheaegselt põhimällu. Seetõttu ei tohi vastusõnumi keha suurus ületada 10 MB (turvaserveri vaikeseadistus).

Vastuses võib olla kuni `batchmaxsize` kirjet. Kui vastuses on kirjeid vähem kui `batchmaxsize`, siis see tähendab, et EHS-l ei ole vastuse koostamise hetkel rohkem andmeid e-hääletamiste kohta.

VIS3 peab suutma töödelda erineva `batchmaxsize` väärtusega vastuseid.

### Töötluse ülevaade

EHS registreerib e-hääletamise fakte. Igale e-hääletamisele omistab EHS järjenumbri. E-hääletamise fakte võib hoida nt järjenumbri järgi indekseeritud tabelis - siis on VIS3-i päringutele vastamine kiire ja efektiivne - kuid see on EHS siseasi. Paki maksimaalsuurus peaks olema EHS seadistuses määratav.

VIS3 saadab valimiste aktiivsel perioodil regulaarselt päringuid EHS-i X-tee teenuse otspunkti "Viimane järjenumber". EHS saadab vastuses viimase e-hääletamise järjenumbri.

Viimase järjenumbri alusel leiab VIS3, kas VIS3-e kantud e-hääletamiste andmed on EHS-iga sünkroonis. Kui viimane järjenumber osutab, et EHS-s on lisandunud kirjeid, mis vajavad VIS3-e kandmist, siis saadab VIS3 järjest päringud EHS X-tee teenuse otspunkti "e-hääletamiste pakk", alates esimesest järjenumbrist, mis on VIS3-s puudu; EHS saadab küsitud paki; VIS3 salvestab saadud andmed ja saadab järgmise paki päringu.

Päringute "e-hääletamiste pakk" töötlusloogika EHS-i poolel peab võimaldama VIS3-l päringuid esitada mistahes järjekorras ja kuitahes palju kordi.

Töötlus peab olema idempotentne (samajõuline) - selles mõttes, et VIS3 võib päritud andmeid igal ajal uuesti küsida. Uuesti pärimisega ei tohi tekkida duubelandmeid, tähendusnihkeid ega kinnijooksmisi.

Tehnilise taustateabena märgime, et EHS hoiab e-hääletamise fakte mitte relatsioonilises andmebaasis, vaid etcd mäluteenuses.

### Kohaletoimetamise garantii
Sõltuvalt e-hääletamise sagedusest ja EHS jõudlusest ning seadistusest võib EHS-s juhtuda, et e-hääletamise väga suure sageduse perioodil e-hääletamise faktile järjenumbri omistamine ajalõpu (ingl timeout) tõttu ebaõnnestub. Selline e-hääletamise fakt jääb e-hääletamiste nimekirjas VIS3-e edastamata.

Seega e-hääletamiste nimekiri ei anna kõigi e-hääletamiste VIS3-e jooksvalt kohaletoimetamise garantiid. VIS3-e kohaletoimetatud fakte tuleb käsitada informatiivsetena. E-hääletamise faktid toimetatakse VIS3-e kindlalt täielikus koosseisus e-hääletanute (lõplikus) nimekirjas, pärast e-hääletamise perioodi lõppu (eraldi liides EHS ja VIS3 vahel).

### Teenuse pakkumise ajaline ulatus

Konkreetse valimissündmuse kohta pakub EHS e-hääletamiste nimekirja ainult piiratud perioodil. See periood hõlmab e-hääletamise perioodi (kehtiva õiguse kohaselt 6 päeva) koos lühikeste siirdeperioodidega enne ja pärast.

Kui EHS valimissündmuse kohta e-hääletamiste nimekirja enam ei paku, siis päringud otspunktidesse "Viimane järjenumber" ja "e-hääletanute pakk" saavad HTTP vastuskoodi `410 Gone`.

Valimissündmuste loetelu pakub EHS pidevalt (otspunkt "Valimissündmuste loetelu"). Loetelus on valimissündmused, mille kohta EHS on valmis e-hääletamiste nimekirja pakkuma.

### Kontroll ja veaolukordade käsitlemine

Eeldatakse, et EHS-i poolt VIS3-le väljastatav on korrektne ja muutumatu (e-hääletamiste nimekirja piires). Parandus- ja muutmiskirjeid käesolev protokoll ei sisalda.

Pärast e-hääletamise lõppu, valimispäeval, edastatakse EHS-st VIS3-le e-hääletanute lõplik nimekiri). See edastus on spetsifitseeritud: [E-hääletanute nimekiri](https://github.com/e-gov/VIS3-EHS/blob/main/4_e_haaletanute_nimekiri/SPEC.md).

VIS3 operaator laeb e-hääletanute lõpliku nimekirja VIS3-e.

Jooksva nimekirja ja lõpliku nimekirja erinevuse korral loetakse "tõe allikaks" lõplik nimekiri.
Seega: 1) kui lõplikus nimekirjas on isik, kes hääletamisperioodil edastatud nimekirjas puudub, siis loetakse isik e-hääletanuks; 2) kui hääletamisaegses nimekirjas on isik, kes lõplikus nimekirjas puudub, siis märge sellise isiku kohta VIS3-s küll säilitatakse, kuid töötluses ja toimingutes lähtutakse lõplikust nimekirjast, s.t loetakse, et isik ei ole e-hääletanud. VIS3 peab omama võimekust hääletamisaegse ja lõpliku nimekirja erinevust avastada ja operaatorile teada anda.

EHS-st VIS3-e edastatud e-hääletamise faktid on VIS3 kasutajaliideses nähtavad valimiste korraldajale, vaates "Valijaga seotud toimingute ajalugu". Kui valija e-hääletas mitu korda, siis on valijaga seotud toimingute ajaloos esitatud kõik valija e-hääletamised. E-hääletamise juures on näha kuupäev ja kellaaeg, millal e-hääletamise fakt EHS-st VIS3-e edastati.

## Kirjandus

\[1] Git. https://git-scm.com/.

\[2] CouchDB Replication Protocol, https://guide.couchdb.org/draft/replication.html.

\[3] Principle and application of Rsync algorithm. https://developpaper.com/principle-and-application-of-rsync-algorithm/.
