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
|   ``String``  |                        "Hey Beavis, I'm a string huh!"                                     |
|   ``Number``  |                              0, 1, 2, 3, 5                                                 |
|   ``Boolean`` |                           yes, no, true, false                                             |
|   ``Alien``   |(Things that requires more explanation in order to express all developer's craziness about) |

All configuration is defined in the ``field = value`` form. Being one configuration (``field-value``) per line.

The main section in a ``cherry configuration file`` is called ``cherry.root`` and until now it only admits one information inside.
This piece of information is the server's hostname. If your server has no name you can use the literal IP address as follows:

```
        cherry.root (
            # This is a comment, sorry I forgot to talk about comments.
            servername = "192.30.70.3"
        )
```

There is a section where you actually open your ``chat rooms``. This section is called ``cherry.rooms``.
There ``alien values`` are needed. This ``alien value`` must be in this form: ``[room_name]:[listen_port]``.
So take a look at the definition sample right below:

        cherry.rooms (
            aliens-on-earth:8810
            foobaroom:8811
            wonkies-lounge:8812
            backyard-science:8813
        )

Each room opened inside ``cherry.rooms`` section features specific sections that must be adjusted in order to be created
at the moment that you run ``Cherry``. The ``Table 2`` summarizes these sections.

**Table 2**: Specific room's sections.

|              **Sections**                      |                    **Used for**                                        |
|:----------------------------------------------:|:----------------------------------------------------------------------:|
|       ``cherry.[room-name].template``          |                  templates definition                                  |
|       ``cherry.[room-name].actions``           |                  actions definition                                    |
|       ``cherry.[room-name].actions.templates`` |                  actions templates definition                          |
|       ``cherry.[room-name].images``            |                  images definition                                     |
|    ``cherry.[room-name].images.url``           |                  images resources definition                           |
|        ``cherry.[room-name].misc``             |                  generic configurations for this room                  |

All information inside ``Table 2`` must be a mess for you. For this reason, firstly, we need to understand some concepts:
``templates``, ``actions``, ``images`` and ``misc configs``.

## What are templates?

Templates can be understood as ``HTML`` data bringing some special makers which are processed (expanded) before sent. When
sent it means sent to the room clients. These special markers on ``Cherry`` are delimited by ``{{.`` and ``}}``. ``Table 3``
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

Actions are the ways how users can communicate each other. Your chat room for example can admit that a user: "talks", "screams" and "mutters".

The way to define it for a room is as follows:

```
        cherry.aliens-on-earth.actions {
            a01 = "talks"
            a02 = "screams"
            a03 = "mutters"
            a04 = "(IGNORE)"
            a05 = "(STOP IGNORE)"
        }

```

Each action definition should be: ``<action-identifier> = <action label string>``.

Depending on action it is possible to format the message in a specific way. For this reason there is another section
called ``cherry.[room-name].actions.templates`` where this must be defined.

```
        cherry.aliens-on-earth.actions {
            a01 = "aliens-on-earth/templates/actions/a01.html"
            a02 = "aliens-on-earth/templates/actions/a02.html"
            a03 = "aliens-on-earth/templates/actions/a03.html"
            a04 = "aliens-on-earth/templates/actions/a04.html"
            a05 = "aliens-on-earth/templates/actions/a05.html"
        }
```

Each action template definition should be: ``<action-identifier-previous-defined-inside-actions> = <string path to a HTML template>``.

## What are images?

Similar to the ``actions`` the ``images`` are labels which the user can choose inside a combo when sending messages. This message
when formatted will include an well-known image. Usually an image should be tematic. Well, things like smiles, emojis, etc.

The images are configurated using two sections. The first one defines the identifiers and their labels.

```
        cherry.aliens-on-earth.images {
            i01 = "glad"
            i02 = "mad"
            i03 = "abducted"
        }
```

Now with the identifiers and labels properly defined it is necessary indicate the URL from each resource (an image in this case).

```
        cherry.aliens-on-earth.images.url {
            i01 = "http://www.nasa.org/chat51/glad.gif"
            i02 = "http://www.nasa.org/chat51/mad.gif"
            i03 = "http://www.nasa.org/chat51/abducted.gif"
        }
```

## What about the misc config?

Misc configurations are generic configurations for a specific room. It can be accessed from section called: ``cherry.[room-name].misc``.
Take a look at the ``Table 4`` in order to see what can be configurated in this section.

|      **Configuration**                   |               **What it does**                           |  **Data type**     |
|:----------------------------------------:|:--------------------------------------------------------:|:------------------:|
|       ``join-message``                   | Defines a message that is displayed when a new user joins|      ``string``    |
|       ``exit-message``                   | Defines a message that is displayed when a user exits    |      ``string``    |
|       ``on-ignore-message``              | Message that confirms an ignore action                   |      ``string``    |
|       ``on-deignore-message``            | Message that confirms a (de)ignore action                |      ``string``    |
|       ``greeting-message``               | Defines a greeting message                               |      ``string``    |
|       ``private-message-maker``          | Defines a string that indicates a private message        |      ``string``    |
|       ``max-users``                      | Defines the maximum of users allowed for this room       |      ``number``    |
|       ``allow-brief``                    | Defines if briefs are allowed or not                     |      ``boolean``   |
|       ``all-users-alias``                | Defines the alias which represents everybody in the room |      ``string``    |
|       ``ignore-action``                  | Defines the action-id used as ignore command             |      ``string``    |
|       ``deignore-action``                | Defines the action-id used as (de)ignore commnad         |      ``string``    |

Follows a definition sample:

```
        cherry.aliens-on-earth.misc (
            join-message = "joined..."
            exit-message = "has left..."
            on-ignore-message = "(only you can see this) is ignoring "
            on-deignore-message = "(only you can see this) is not ignoring "
            greeting-message = "welcome"
            private-message-marker = "(private)"
            max-users = 10
            allow-brief = yes
            all-users-alias = "EVERYBODY"
            ignore-action = "a04"
            deignore-action = "a05"
        )
```

## Some tricks

It is not a good practice define the entire configuration in just one file. The ``Cherry`` configuration's language implements
support for code "importation". The way to do it is:

```
        #
        # "cherry.config"
        # Description: the main configuration file.
        #

        cherry.branch aliens-on-earth.config
        cherry.branch backyard-science.config
        cherry.branch wonkies-lounge.config
```

Congrats, now your ``Cherry`` tree has branches! :) Cut off one branch from it is pretty simple, just comment it.
In the sample case above as effect a room will stop being created.

## Opening your first chat room

I know is rather confuse read this kind of descriptions without any concrete example. From now on we will compose each
document necessary to create a chat room.

On ``Cherry`` there are 3 kinds of documents (HTML documents):

1. The join form (where the user chooses his nickname and color for it).
2. The chat room (with 3 parts: top, body and banner).
3. The post-exit ``HTML`` document.

### Adding the chat room brief support to your server

### Adding user find support to your server
