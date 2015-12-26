# Cherry doc

``Cherry`` is an application that works through a configuration file. This configuration file has his own way to express things.
This configuration divide stuffs in sections. A section can be understood as:

```
        [section-name] (
            [all configuration stuff goes here]
        )
```

Configurations can have some data types such as strings, alien values, number and booleans. Take a look at the ``Table 1``
to see samples from these data types.

**Table 1**: Data types sample.

| **Data Type** |                                         **Sample**                                         |
|:-------------:|:------------------------------------------------------------------------------------------:|
|   ``String``  |                        "Hey Beavis, I'm an string huh!"                                    |
|   ``Number``  |                              0, 1, 2, 3, 5                                                 |
|   ``Boolean`` |                           yes, no, true, false                                             |
|   ``Alien``   |(Things that requires more explanation in order to express all developer's craziness about) |

All configuration is defined in the ``field = value`` form. Being one configuration (``field-value``) per line.

The main section in a ``cherry configuration file`` is called ``cherry.root`` and until now this only admits one information inside.
This piece of information is the server's hostname. If your server does not have a name you can use the literal IP address as
follows:

        cherry.root (
            # This is a comment, sorry I forgot to talk about comments.
            servername = "192.30.70.3"
        )

There is a section where you actually open your ``chat rooms``. This section is called ``cherry.rooms``.
There ``alien values`` are needed. This alien value must be in this form: ``[room_name]:[listen_port]``.
So take a look at the definition sample right below:

        cherry.rooms (
            aliens-on-earth:8810
            foobaroom:8811
            wonkies-lounge:8812
            backyard-science:8813
        )

Each room opened inside ``cherry.rooms`` section features specific sections that must be adjusted in order to be created
in the moment that you run ``Cherry``. The ``Table 2`` summarizes these sections.

**Table 2**: Specific room's sections.

|              **Sections**                      |                    **Used for**                                        |
|:----------------------------------------------:|:----------------------------------------------------------------------:|
|       ``cherry.[room-name].template``          |                  templates definition                                  |
|       ``cherry.[room-name].actions``           |                  actions definition                                    |
|       ``cherry.[room-name].actions.templates`` |                  actions templates definition                          |
|       ``cherry.[room-name].images``            |                  images definition                                     |
|    ``cherry.[room-name].images.url``           |                  images resources definition                           |
|        ``cherry.[room-name].misc``             |                  generic configurations for this room                  |

All information inside ``Table 2`` must be a confusion for you. For this reason we need before understand some concepts:
``templates``, ``actions``, ``images`` and ``misc configs``.

## What are templates?

Templates can be understood as ``HTML`` data bringing some special makers which are processed (expanded) before sent. When
sent it means sent to the clients. These special markers in ``Cherry`` are delimited by ``{{.`` and ``}}``. ``Table 3``
summarizes each special marker supported until now.

**Table 3**: Special markers.

|               **Marker**                       |                         **Expanded for**                               |
|:----------------------------------------------:|:----------------------------------------------------------------------:|
|            ``{{.nickname}}``                   |                      The user's nickname                               |
|            ``{{.session-id}}``                 |                      The user's session-id                             |
|            ``{{.color}}``                      |                      The user nickname's color code                    |
|            ``{{.ignorelist}}``                 |                      Users that one user are ignoring                  |
|            ``{{.hour}}``                       |                      The current server hour                           |
|            ``{{.minute}}``                     |                      The current server minute                         |
|            ``{{.second}}``                     |                      The current server second                         |
|            ``{{.greeting-message}}``           |                      The configurated greeting message                 |
|            ``{{.join-message}}``               |                      The configurated join message                     |
|            ``{{.exit-message}}``               |                      The configurated exit message                     |
|            ``{{.on-ignore-message}}``          |                      The configurated ignore message                   |
|            ``{{.on-deignore-message}}``        |                      The configurated "(de)ignore" message             |
|            ``{{.max-users}}``                  |                      The maximium users supported by this room         |
|            ``{{.all-users-alias}}``            |                      Alias that represents everybody (broadcast)       |
|            ``{{.action-list}}``                |                      Action list to be included in the "talk-banner"   |
|            ``{{.image-list}}``                 |                      Image list to be included in the "talk-banner"    |
|            ``{{.users-list}}``                 |                      Users list to be included in the "talk-banner"    |
|            ``{{.top-template}}``               |                      The Top template                                  |
|            ``{{.body-template}}``              |                      The body template                                 |
|            ``{{.banner-template}}``            |                      The banner template                               |
|          ``{{.highlight-template}}``           |                      The highlight which diffs personal messages       |
|          ``{{.entrace-template}}``             |                      The room's entrace form                           |
|          ``{{.exit-template}}``                |                      The room's post-exit document                     |
|          ``{{.nickclash-template}}``           |                      The room's nickclash warning document             |
|          ``{{.last-public-messages}}``         |                      Public messages that can be used to compose briefs|
|          ``{{.servername}}``                   |                      The configurated server name                      |
|          ``{{.listen-port}}``                  |                      The room's listen port                            |
|          ``{{.room-name}}``                    |                      The room's name                                   |
|          ``{{.users-total}}``                  |                      The current amount of connected users on that room|
|          ``{{.message-action-label}}``         |                      The label from a choosen action                   |
|          ``{{.message-whoto}}``                |                      The message destination user                      |
|          ``{{.message-user}}``                 |                      The message source user                           |
|          ``{{.message-colored-user}}``         |                      The message source user (formatted with the color)|
|          ``{{.message-says}}``                 |                      The message data                                  |
|          ``{{.message-image}}``                |                      The message image icon (if this has one)          |
|          ``{{.message-private-marker}}``       |                      The private marker of a private message           |
|          ``{{.brief-last-public-messages}}``   |                      The last public messages (well formatted)         |
|          ``{{.brief-who-are-talking}}``        |                      The user list (well formatted)                    |
|          ``{{.brief-users-total}}``            |                      The users total (well formatted)                  |
|          ``{{.find-result-user}}``             |                      The find result (user nickname)                   |
|          ``{{.find-result-room-name}}``        |                      The find result (user room)                       |
|          ``{{.find-result-users-total}}``      |                      The find result (total of users in the user room) |

## What are actions?

## What are images?

## What about the misc config?

