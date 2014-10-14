Backuper - (самый) простой инструмент для backup
========
[![Build Status](https://drone.io/github.com/headmade/backuper/status.png)](https://drone.io/github.com/headmade/backuper/latest)
Backuper server agent


## Фишки:
1. [простая установка](#Установка)
2. [очень простая настройка](#config)
2. [конфигурирование через web-interface](#control)
3. [command-line инструменты]
4. возможность автономной работы
5. хранение backup на Amazon S3
6. базовый функционал - бесплатен навсегда
7. open source агент

# Как пользоваться

Далее следует описание действий, необходимых для настройки резервного копирования одного сервера.

## Регистрация пользователя

Необходимо зарегистрироваться в системе
http://backuper.headmade.pro
(требуется только email и пароль; на этот email будут присылаться только отчёты о backup, мы его никому не продадим).


## Регистрация сервера (backup которого нужен)

Зайти по ссылке и ... TODO ...

При регистрации этому серверу назначается уникальный идентификатор `UUID`, о чём радостно сообщает web-interface, и этот идентификатор вскоре понадобится.


## Установка backup-агента

Для установки backup-агента на сервере, которому необходимо резервное копирование, необходимо запустить:

### Debian

```
sudo sh -c 'echo "deb http://apt.backuper.headmade.pro $(lsb_release -cs) main" >> /etc/apt/sources.list'
wget --quiet -O - http://apt.backuper.headmade.pro/B4C2B02A.asc | sudo apt-key add -
sudo apt-get update
sudo apt-get install backuper
```


### MacOS

`brew install backuper`

### Из исходников

1. Установить OpenSSL
   
   `sudo apt-get install openssl`

1. [install go](https://golang.org/dl/)
2. `git clone https://github.com/headmade/backuper.git && cd backuper && make install`

## Начальная настройка агента

1. `backuper init UUID`

    где `UUID` - UUID, назначенный при регистрации этого сервера.
    
    В результате этой команды агент идентифицируется в web-interface,
    и отныне можно конфигурировать backup-настройки через web.

2. Для хранения резервных копий на Amazon S3 запустите:

    `backuper provider aws <AWS_ACCESS_KEY_ID>  <AWS_SECRET_ACCESS_KEY>`
    
    где `AWS_ACCESS_KEY_ID` и `AWS_SECRET_ACCESS_KEY` - ключи доступа к S3.

3. Настройка шифрования:

    `backuper provider encrypt <SOME_STRONG_PASSWORD>`
    
    где `SOME_STRONG_PASSWORD` будеи использоваться для шифрования backup-файлов. Это значение не выходит за пределы резервируемого сервера.



## Настройка агента

Агент можно сконфигурировать полностью через web-interface.
Пытливые и любознательные могут сделать это вручную, исправляя конфигурационный json-файл, но в абсолютном большинстве случаев можно обойтись без этого.

1. в Web-interface, выбрать свежесозданный сервер и зайти в настройки backup
![](http://puu.sh/c34dA/919f4f322e.png)

2. Нажать кнопку "Настроить..."

3. Настроить периодичность выполнения backup.

    ![](http://puu.sh/c36he/43ca5f5601.png)

4. Настроить backup-таски

    На данный момент сервис умеет бакапить 2 вида данных:
    - локальный файл/папку
    - базу postgres
    
    Для каждого вида доступны свои, интуитивно-очевидные настройки.
    
    Чтобы сделать резервную копию нескольких баз данных или нескольких папок, создайте несколько backup-тасков нужного вида.

5. Настроить место хранения backup-файлов

    На данный момент поддерживаются следующие варианты хранения:
    ![](http://puu.sh/c35sz/27439f0e45.png)
    Перед отправкой backup-файл шифруется известным только вам паролем.
    Незашифрованные данные никуда не передаются и нигде не хранятся за пределами резервируемого серера.

6. Укажите папку для хранения временных файлов (например, дампов баз данных).

    ![](http://puu.sh/c36nB/af2c44a43e.png)

7. И нажмите "Сохранить"

Начиная со следующей минуты серверного времени, при наступлении указанного в настройках момента будет выполняться backup.

Всё :)



## Contributing

1. Fork it
2. Create your feature branch (```git checkout -b my-new-feature```).
3. Commit your changes (```git commit -am 'Added some feature'```)
4. Push to the branch (```git push origin my-new-feature```)
5. Create new Pull Request
