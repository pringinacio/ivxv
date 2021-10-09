..  IVXV kogumisteenuse haldusjuhend

Lisad
=====

.. _utiliidid:

Utiliidid
---------

Kogumisteenuse haldamise käsureautiliitide ülevaade ja abiteave.

.. contents:: .
   :local:
   :depth: 1


Andmehoidla utiliidid
^^^^^^^^^^^^^^^^^^^^^

.. include:: utiliitide-abiteave/ivxv-create-data-dirs.inc

.. include:: utiliitide-abiteave/ivxv-db-reset.inc

.. include:: utiliitide-abiteave/ivxv-db-dump.inc


Teenuse seisundi utiliidid
^^^^^^^^^^^^^^^^^^^^^^^^^^

.. include:: utiliitide-abiteave/ivxv-status.inc

.. include:: utiliitide-abiteave/ivxv-service.inc


Sündmuste logi utiliidid
^^^^^^^^^^^^^^^^^^^^^^^^

.. include:: utiliitide-abiteave/ivxv-eventlog-dump.inc


Kasutajate halduse utiliidid
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. include:: utiliitide-abiteave/ivxv-users-list.inc


Seadistusutiliidid
^^^^^^^^^^^^^^^^^^

.. include:: utiliitide-abiteave/ivxv-collector-init.inc

.. include:: utiliitide-abiteave/ivxv-cmd-load.inc

.. include:: utiliitide-abiteave/ivxv-config-validate.inc

.. include:: utiliitide-abiteave/ivxv-config-apply.inc

Seadistuste rakendamine hallatavatele teenustele on võimalik siis, kui
haldusteenusesse on laaditud kogumisteenuse tehnilised seadistused.

Seadistuste rakendamise järjekord:

#. Tehnilised seadistused koos usaldusjuure seadistustega.

   #. Teenuse tarkvara paigaldamine;

   #. Haldusteenuse ligipääsu loomine hallatava teenuse kontole;

   #. Teenuse logimisseadistuste rakendamine;

   #. Haldusteenuse ligipääsu eemaldamine teenuse hosti juurkasutaja kontole
      (ainult juhul, kui teenusmasinas pole rohkem seadistamata teenuseid);

   #. Usaldusjuure rakendamine teenusele;

   #. Tehniliste seadistuste rakendamine teenusele;

#. Valikute nimekiri;

#. Ringkondade nimekiri;

#. Valijate nimekirjad;

Logikogumisteenus erineb teistest hallatavatest teenustest:

#. Logikogumisteenus seadistatakse enne teisi teenuseid, et tagada võimalikult
   varajane logi kogumine.

#. Logikogumisteenustele ei rakendata muid seadistusi peale logikogumisteenuse
   seadistuste (usaldusjuure seadistusi, kogumisteenuse tehnilised seadistusi
   ja valimiste seadistusi logikogumisteenus ei vaja).

Valimisnimekirjade (valikute ja valijate nimekirjad) rakendamine tähendab
nimekirja ülekandmist talletusteenusesse vastavat nimekirja teenindava teenuse kaudu.

Näiteks valikute nimekiri rakendatakse vaid ühele (juhuslikult valitud)
nimekirjateenusele, mis kannab nimekirja talletusteenusesse. Talletusteenuse
kaudu on nimekiri kättesaadav kõigile teistele nimekirjateenustele.

.. include:: utiliitide-abiteave/ivxv-voter-list-download.inc

.. include:: utiliitide-abiteave/ivxv-secret-load.inc

.. include:: utiliitide-abiteave/ivxv-copy-log-to-logmon.inc

.. include:: utiliitide-abiteave/ivxv-update-packages.inc

.. include:: utiliitide-abiteave/ivxv-backup-crontab.inc


Andmete eksportimise ja varundamise utiliidid
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. include:: utiliitide-abiteave/ivxv-export-votes.inc

.. include:: utiliitide-abiteave/ivxv-backup.inc

.. include:: utiliitide-abiteave/ivxv-generate-processor-input.inc

.. include:: utiliitide-abiteave/ivxv-voterstats.inc

.. include:: utiliitide-abiteave/ivxv-voting-sessions.inc


Deemonid
^^^^^^^^

.. include:: utiliitide-abiteave/ivxv-agent-daemon.inc


Sisemised utiliidid
^^^^^^^^^^^^^^^^^^^

.. attention::

   Sisemised utiliidid on kasutusel haldusdeemoni poolt alamteenuste
   haldamiseks ja neid ei ole reeglina tarvis eraldi käivitada.

.. include:: utiliitide-abiteave/ivxv-admin-helper.inc

.. include:: utiliitide-abiteave/ivxv-admin-sudo.inc


Seadistusfailid
---------------

.. _ivxv-logcollector.conf:

Logikogumisteenuse seadistusfail
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. literalinclude:: ../../../common/collector/config/rsyslog-logcollector.conf
   :name: /etc/rsyslog.d/ivxv-logcollector.conf
   :language: text
   :linenos:


Lisaseadistused
---------------

.. _configure-ssh-idcard-auth:

SSH kasutajate autentimine ID-kaardi abil
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

SSH-teenusesse on võimalik autentida ID-kaardi avaliku võtmega abil, kasutades
selleks PKCS#11 toega SSH-klienti kitty.exe (http://kitty.9bis.net/).

Turvakaalutustel tuleks keelata haldusliidese SSH-teenusesse parooliga
autentimine. Parooliga autentimise keelamiseks tuleb seadistusfailis
:file:`/etc/ssh/sshd_config` määrata parameetri ``PasswordAuthentication``
väärtuseks ``no``::

   # To disable tunneled clear text passwords, change to no here!
   PasswordAuthentication no

Volitatud kasutajate faili asukoht (:file:`/etc/ssh/kasutajad`) tuleb failis
:file:`/etc/ssh/sshd_config` määrata parameetriga ``AuthorizedKeysFile``:

   AuthorizedKeysFile /etc/ssh/kasutajad

.. important::

   Seadistusfailis /etc/ssh/sshd_config tehtud muutuse rakendamiseks tuleb
   SSH teenus taaskäivitada::

      # service ssh restart
      [ ok ] Restarting OpenBSD Secure Shell server: sshd.

ID-kaardi isikutuvastamise sertifikaadiga autenditava kasutaja ülesseadmine
käib järgmiselt:

#. Kasutajale konto loomine:

   .. code-block:: shell-session

      # adduser --disabled-password kasutajanimi
      # usermod -a -G www-data kasutajanimi

#. Kasutaja ID-kaardi isikutuvastamise sertifikaadi salvestamine PEM-vormingus
   faili :file:`usercert.cer` (ID-kaardi haldusvahendi abil);

#. Sertifikaadist kasutaja avaliku võtme eraldamine ja salvestamine faili
   :file:`userpubkey.pem`:

   .. code-block:: shell-session

      # openssl x509 -in usercert.cer -pubkey -noout > userpubkey.pem

#. Avaliku võtme teisendamine PKCS#8 vormingusse, kasutaja tunnusega
   varustamine ja salvestamine SSH volitatud kasutajate faili
   :file:`/etc/ssh/kasutajad`:

   .. code-block:: shell-session

      # KEY=$(ssh-keygen -i -m PKCS8 -f userpubkey.pem)
      # echo "$KEY kasutaja@eesti.ee" >> /etc/ssh/kasutajad

#. Kontrollimine, kas lisatud kirje on kujul ``ssh-rsa PKCS8-võti``
   kasutajatunnus:

   .. code-block:: shell-session

      # tail -1 /etc/ssh/kasutajad
      ssh-rsa AAAAB3NzaC1yc2EAAAAELGuiTwAAAIEAxZf/TuSrGJEU1PlfkY9jJ33VOYVZ9Vao0Uiytlf8
      7HJu/78fCIB7m05J7ibpMhsZoZ4DElU7ve0VwbvdDS3srh1OhiQcUjpznTlx4rIM1vkHwadrHtmF+BNi
      DwbLbbdD5y3puGcLH+sLuwba6Vuc3aU0QuqzenYmY9pV7w9y0wc= kasutaja@eesti.ee


Andmehoidla
-----------

Haldusteenuse andmeid hoitakse failisüsteemis ja andmebaasis. Failisüsteemis
hoitakse andmeid, mis on pärit välistest süsteemidest ja on haldusteenusesse
üle kantud faili kujul. Andmebaasis hoitakse andmeid, mis on genereeritud
haldusteenuse töö käigus.


Failisüsteemis hoitavad andmed
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

*  :file:`/etc/ivxv/` -- kogumisteenusele rakendatud ja hetkel kehtivad
   seadistus- ja nimekirjafailid;

*  :file:`/var/lib/ivxv/` -- kogumisteenuse haldusteenuse andmefailid;

*  :file:`/var/lib/ivxv/admin-ui-data/` -- haldusteenuse veebiliidese jaoks
   serveeritavad JSON-failid;

*  :file:`/var/lib/ivxv/admin-ui-data/status.json` -- kogumisteenuse seisundi
   koondandmed;

*  :file:`/var/lib/ivxv/admin-ui-permissions/` -- haldusteenuse veebiliidese
   kasutajaõigused (Apache veebiserveri jaoks);

*  :file:`/var/lib/ivxv/ballot-box/` -- allalaaditava e-valimiskasti salvestamise kataloog;

*  :file:`/var/lib/ivxv/commands/` -- kogumisteenuse juhtimiseks rakendatud
   korraldusfailide ajalugu;

*  :file:`/var/lib/ivxv/commands/<command-type>-<timestamp>.bdoc` --
   digitaalselt allkirjastatud korraldus ASiC-E vormingus.

*  :file:`/var/lib/ivxv/commands/<command-type>-<timestamp>.json` --
   korralduse olekufail JSON-vormingus.

*  :file:`/var/lib/ivxv/db/` -- haldusteenuse andmebaasi kataloog;

*  :file:`/var/lib/ivxv/db/ivxv-management.db` -- haldusteenuse andmebaasi
   fail;

*  :file:`/var/lib/ivxv/ivxv-management-events.log` -- haldusteenuse sündmuste
   logi;

*  :file:`/var/lib/ivxv/service/` -- muud teenusespetsiifilised failid
   (nt. registreerimisvõtmest eraldatud avalik võti);

*  :file:`/var/lib/ivxv/upload/` -- kogumisteenusesse veebiliidese kaudu
   laaditud failid;

Andmebaasis hoitavad andmed
^^^^^^^^^^^^^^^^^^^^^^^^^^^

Andmevälja nimi ja kirjeldus:

* ``collector/state`` -- kogumisteenuse olek;

* ``config/election`` -- kogumisteenuses rakendatud valimiste seadistusele
  digiallkirja andnud volitatud kasutaja andmed kujul ``<CN> <timestamp>``;

* ``config/technical`` -- kogumisteenuses rakendatud tehnilisele seadistusele
  digiallkirja andnud volitatud kasutaja andmed kujul ``<CN> <timestamp>``;

* ``config/trust`` -- kogumisteenuses rakendatud usaldusjuure seadistusele
  digiallkirja andnud volitatud kasutaja andmed kujul ``<CN> <timestamp>``;

* ``election/election-id`` -- valimiste identifikaator;

* ``election/electionstart`` -- valimiste algusaeg;

* ``election/electionstop`` -- valimiste lõpuaeg;

* ``election/servicestart`` -- kogumisteenuse käivitamise aeg;

* ``election/servicestop`` -- kogumisteenuse seiskamise aeg;

* ``host/<hostname>/state`` -- teenushosti seisund;

* ``list/choices`` -- haldusteenusesse laaditud valikute nimekirjale
  digiallkirja andnud volitatud kasutaja andmed kujul ``<CN> <timestamp>``;

* ``list/choices-loaded`` -- nimekirjateenustesse laaditud valikute
  nimekirjale digiallkirja andnud volitatud kasutaja andmed kujul ``<CN>
  <timestamp>``;

* ``list/districts`` -- nimekirjateenustesse laaditud ringkondade
  nimekirjale digiallkirja andnud volitatud kasutaja andmed kujul ``<CN>
  <timestamp>``;

* ``list/districts-loaded`` -- nimekirjateenustesse laaditud ringkondade
  nimekirjale digiallkirja andnud volitatud kasutaja andmed kujul ``<CN>
  <timestamp>``;

* ``list/voters0000`` -- haldusteenusesse laaditud valijate algnimekirjale
  digiallkirja andnud volitatud kasutaja andmed kujul ``<CN> <timestamp>``;

* ``list/voters<list-number>`` (``list-number >= 01``) -- haldusteenusesse
  laaditud valijate muudatusnimekirja allalaadimise hetke ajatempel;

* ``list/voters<list-number>-state`` -- nimekirjateenustesse laaditud valijate
  nimekirja olek.

  Võimalikud väärtused:

  1. ``PENDING`` - laaditud haldusteenusesse;

  2. ``APPLIED`` - rakendatud nimekirjateenusele;

  3. ``INVALID`` - nimekiri on märgitud vigaseks ja ootab halduri otsust
     vahelejätmise kohta (ainult muudatusnimekirja korral);

  4. ``SKIPPED`` - nimekiri on vahele jäetud (ainult muudatusnimekirja korral).

* ``logmonitor/address`` -- seireteenuse aadress või võrgunimi;

* ``logmonitor/last-data`` -- viimase seireteenusest statistikafaili hankimise
  aeg;

* ``user/<idcode>`` -- haldusteenuse kasutaja nimi ja rollid kujul
  ``<surname,name> <role>[,<role>]``;

* ``service/<service-id>/service-type`` -- Teenuse liik;

* ``service/<service-id>/technical-conf-version`` -- Teenusele rakendatud
  tehnilise seadistuse versioon;

* ``service/<service-id>/election-conf-version`` -- Teenusele rakendatud
  valimiste seadistuse versioon;

* ``service/<service-id>/network`` -- Teenusele alamvõrgu nimi;

* ``service/<service-id>/state`` -- Teenuse olek;

* ``service/<service-id>/ping-errors`` -- Teenuse elusoleku kontrollimise
  järjestikuste vigade arv;

* ``service/<service-id>/last-data`` -- Teenuse viimase oleku hankimise aeg;

* ``service/<service-id>/ip-address`` -- Teenuse IP-aadress;

* ``service/<service-id>/bg_info`` -- Teenuse taustainfo stringina (näiteks
  elusoleku kontrolli käigus genereeritud veateade);

* ``service/<service-id>/backup-times`` -- Varundusteenuse automaatvarunduse
  kellaajad;

* ``service/<service-id>/mid-token-key`` -- Mobiil-ID tugiteenuse
  identsustõendi võtmefaili kontrollsumma (SHA256);

* ``service/<service-id>/tls-cert`` -- Teenuse TLS-sertifikaadi faili
  kontrollsumma (SHA256);

* ``service/<service-id>/tls-key`` -- Teenuse TLS-sertifikaadi võtmefaili
  kontrollsumma (SHA256);

* ``service/<service-id>/tspreg-key`` -- Hääletamisteensue ajatempliteenuse
  signeerimisvõtme faili kontrollsumma (SHA256);

Kasutatud tähised:

* ``<command-type>`` -- korralduse liik:

   #. ``trust`` -- usaldusjuure seadistused;

   #. ``technical`` kogumisteenuse seadistused;

   #. ``election`` valimiste seadistused;

* ``<CN>`` -- ID-kaardi CN väli kujul ``PEREKONNANIMI,EESNIMI,ISIKUKOOD``;

* ``<config-type>`` on seadistuse liik. Usaldusjuure seadistus on ``trust``,
  valimiste seadistus on ``election`` ja kogumisteenuse tehniline
  seadistus on ``tech``;

* ``<hostname>`` teenushosti nimi;

* ``<list-number>`` valimisnimekirja kahekohaline järjekorranumber, esimene nimekiri
  kannab numbrit 01.

* ``<service-id>`` teenuse identifikaator kogumisteenuse seadistustest;

* ``<timestamp>`` on ajatempel ISO-8601 vormingus.


.. _etcd-zabbix:

Klastri seisundi monitoorimine Zabbixiga
----------------------------------------

Etcd klaster tagab süsteemi toimimise ka olukorras, kus mõni klastriliige kaotab
töövõime (krahh, võrguühenduse kadumine jms.). Siiski on oluline selliseid
sündmuseid monitoorida ning nende algpõhjus tuvastada. Etcd krahhimise
tuvastamiseks tuleb talletusteenuste logidest (``ivxv-YYYY-MM-DD.log``) monitoorida
``ivxv.ee/service/storage.EtcdTerminatedError`` kirjet.

Täiendavalt saab etcd käsureakliendiga küsida klastri liikmete olekut. Kuna
IVXV klastris on kõik klient-päringud autenditud, siis tuleb korraldus
käivitada mõnes ``ivxv-storage`` teenuse masinas kasutajakonto
``ivxv-storage`` (või juurkasutaja) õigustes:

.. code-block:: shell-session

   # ivxv-storage@ivxv1:~$ env ETCDCTL_API=3 etcdctl \
         --cacert /var/lib/ivxv/service/storage@storage1.ivxv.ee/ca.pem \
         --cert /var/lib/ivxv/service/storage@storage1.ivxv.ee/tls.pem \
         --key /var/lib/ivxv/service/storage@storage1.ivxv.ee/tls.key \
         --endpoints ivxv1:2379,ivxv2:2379,ivxv3:2379 \
         endpoint status

   ivxv1:2379, 2d0df029f29770a4, 3.2.17, 25 kB, true, 12, 15
   ivxv2:2379, d4a9ae16c8557764, 3.2.17, 25 kB, false, 12, 15
   ivxv3:2379, e8914f4e0b89b80f, 3.2.17, 25 kB, false, 12, 15

Vastuses on veergude tähendused järgmised:

 #. klastri liige;
 #. klastri liikme identifikaator;
 #. etcd versioon;
 #. baasi suurus (max 8GB ehk 8589934592);
 #. kas konkreetne klastri liige on hetkel juht;
 #. RAFT ametiaeg (sisuliselt toimunud juhi-valimiste arv);
 #. RAFT indeks - etcd kirjutamisoperatsioonide arv (sh.
    konfiguratsiooni muutused).

Monitooringule on oluline parameeter RAFT ametiaeg. Selle väärtuse muutumine
tähendab juhivahetust, mis üldjuhul on seotud probleemidega klastri töös -
olemasolev juht ei vasta piisavalt kiiresti klastri liikmete päringutele.

Käsurea seletus:

 * ``env ETCDCTL_API=3``: kasutame etcd API versiooni 3 (Ubuntu versioonis
   20.04 LTS on ``etcdctl`` API vaikeversioon veel 2);
 * ``--cacert``: usaldame ainult servereid, mille sertifikaat on antud selle CA
   poolt;
 * ``--cert`` ja ``--key``: kasutame klient-autentimiseks ivxv1 talletusteenuse
   sertifikaati ja võtit;
 * ``--endpoints``: millistele serveritele päring saata. Siin võib kõigi kolme
   asemel ka ainult ühe loetleda: sellisel juhul on väljundis vaid üks rida.
   Kasulik nt kui Zabbix tahab igas talletusteenuses küsida ainult selle isendi
   kohta;
 * ``endpoint status``: küsime loetletud serverite olekut.


Väljundit on võimalik küsida ka masinloetavas JSON-vormingus
(parameeter ``-w json``):

   .. code-block:: shell-session

      ivxv-storage@ivxv1:~$ env ETCDCTL_API=3 etcdctl \
          --cacert /var/lib/ivxv/service/storage@storage1.ivxv.ee/ca.pem \
          --cert /var/lib/ivxv/service/storage@storage1.ivxv.ee/tls.pem \
          --key /var/lib/ivxv/service/storage@storage1.ivxv.ee/tls.key \
          --endpoints ivxv1:2379,ivxv2:2379,ivxv3:2379 \
          endpoint status -w json

      [{"Endpoint":"ivxv1:2379","Status":{"header":{"cluster_id":1867986262344190226,"member_id":3246514969358332068,"revision":1,"raft_term":12},"version":"3.2.17","dbSize":24576,"leader":3246514969358332068,"raftIndex":15,"raftTerm":12}},
      {"Endpoint":"ivxv2:2379","Status":{"header":{"cluster_id":1867986262344190226,"member_id":15323970619978381156,"revision":1,"raft_term":12},"version":"3.2.17","dbSize":24576,"leader":3246514969358332068,"raftIndex":15,"raftTerm":12}},
      {"Endpoint":"ivxv3:2379","Status":{"header":{"cluster_id":1867986262344190226,"member_id":16758262885041944591,"revision":1,"raft_term":12},"version":"3.2.17","dbSize":24576,"leader":3246514969358332068,"raftIndex":15,"raftTerm":12}}]
