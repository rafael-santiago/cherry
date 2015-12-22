# Cherry doc

``Cherry`` is an application that works through a configuration file. This configuration file has his own way to express things. The idea of this configuration is
divide stuffs in sections. A section can be understood as:

        [section-name] (
            [all configuration stuff goes here]
        )

Configurations can have some data types such as strings, alien values, number and booleans. Look the ``Table 1`` to see samples from these data types.

**Table 1**: Data types sample.

| **Data Type** |                                         **Sample**                                         |
|:-------------:|:------------------------------------------------------------------------------------------:|
|   ``String``  |                        "Hey Beavis, I'm an string huh!"                                    |
|   ``Number``  |                              0, 1, 2, 3, 5                                                 |
|   ``Boolean`` |                           yes, no, true, false                                             |
|   ``Alien``   |(Things that requires more explanation in order to express all developer's craziness about) |

All configuration is defined in the ``field = value`` form. Being one configuration (``field-value``) per line.

The main section in a ``cherry configuration file`` is called ``cherry.root`` and until now this only admits one information inside. This piece of information is the server's hostname.
If you server does not have a name you can use the literal IP address as follows:

        cherry.root (
            # This is a comment, sorry I forgot to talk about comments!
            servername = "192.30.70.3"
        )

There is a section where you actually open your chat rooms. This section is called ``cherry.rooms``. This section needs an aliean value.
This alien value must be in this form: ``[room_name]:[listen_port]``. So take a look at the definition sample right below:

        cherry.rooms (
            aliens-on-earth:8810
            foobaroom:8811
            wonkies-lounge:8812
            backyard-science:8813
        )

Each room opened inside ``cherry.rooms`` section features specific sections that must be configurated in order to be created in the moment that you run ``Cherry``. The ``Table 2``
summarizes these sections.

**Table 2**: Specific room's sections.

|              **Sections**                      |                    **Used for**                                        |
|:----------------------------------------------:|:----------------------------------------------------------------------:|
|       ``cherry.[room-name].template``          |                  templates definition                                  |
|       ``cherry.[room-name].actions``           |                  actions definition                                    |
|       ``cherry.[room-name].actions.templates`` |                  actions templates definition                          |
|       ``cherry.[room-name].images``            |                  images definition                                     |
|    ``cherry.[room-name].images.url``           |                  images resources definition                           |
|        ``cherry.[room-name].misc``             |                  generic configurations for this room                  |

All information inside ``Table 2`` must be a confusion for you. For this reason we need before understand some concepts: templates, actions, images and misc configs.

## What are templates?

## What are actions?

## What are images?

## What about the misc config?

