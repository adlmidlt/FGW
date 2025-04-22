-- Таблица справочников [handbook_id]:
--  0 - конструкторское наименование продукции
--  1 - действия над объектами учета
--  2 - действия над этикеткой
--  3 - цвет продукции
--  4 - принтеры
--  5 - действия для заявок
--  6 - приоритеты
--  7 - статусы заявок
--  8 - компьютеры
--  9 - участки упаковки
-- 10 - участки хранения
-- 11 - объекты учёта
-- 12 - назначение при списании
-- 13 - комментарии к пакет-поддонам
-- 14 - размеры этикеток
-- 15 - типы документов
CREATE TABLE dbo.catalog
(
    id_catalog               INT IDENTITY PRIMARY KEY,           -- генерирует и итерирует serial
    parent_id                INT            DEFAULT 0  NOT NULL, -- для [id_catalog] родительской записи
    handbook_id              INT            DEFAULT 0  NOT NULL, -- номер справочника
    record_index             INT            DEFAULT 0  NOT NULL, -- индекс записи (может повторяться)
    name                     VARCHAR(255)   DEFAULT '' NOT NULL, -- название
    comment                  VARCHAR(5000)  DEFAULT '' NOT NULL, -- комментарий
    handbook_value_int_1     INT            DEFAULT 0  NOT NULL, -- для [handbook_id] (0 - срок годности в месяцах, 1 - объект учета, 3, 4, 9, 10) дополнительное поле 1 справочника (int)
    handbook_value_int_2     INT            DEFAULT 0  NOT NULL, -- дополнительное поле 2 справочника (int)
    handbook_value_decimal_1 DECIMAL(10, 2) DEFAULT 0  NOT NULL, -- для [handbook_id] (10 - возможный процент использования) дополнительное поле 1 справочника (decimal)
    handbook_value_decimal_2 DECIMAL(10, 2) DEFAULT 0  NOT NULL, -- для [handbook_id] (10 - вместимость, S, V) дополнительное поле 2 справочника (decimal)
    handbook_value_bool_1    BIT            DEFAULT 0  NOT NULL, -- для [handbook_id] (9 - переупаковка да/нет, 10 - наличие ЖД путей да/нет9 - переупаковка да/нет, 10 - наличие ЖД путей да/нет) дополнительное поле 1 справочника (boolean)
    handbook_value_bool_2    BIT            DEFAULT 0  NOT NULL, -- дополнительное поле 2 справочника (boolean)
    is_archive               BIT            DEFAULT 0  NOT NULL, -- архивная запись да/нет
    owner_user               UNIQUEIDENTIFIER          NOT NULL, -- uuid владельца записи
    owner_user_datetime      DATETIME                  NOT NULL, -- дата и время записи владельца
    last_user                UNIQUEIDENTIFIER          NOT NULL, -- uuid последнего
    last_user_datetime       DATETIME                  NOT NULL  -- дата и время последней модификации
);

-- Таблица справочников.
CREATE TABLE dbo.handbook
(
    id_handbook INT IDENTITY PRIMARY KEY, -- ИД справочника.
    name        VARCHAR(150) NOT NULL     -- наименование справочника.
);

CREATE PROCEDURE fgw_insert_handbook
    AS
BEGIN
    SET IDENTITY_INSERT handbook ON;
INSERT INTO handbook (id_handbook, name)
VALUES (0, N'Конструкторское наименование'),
       (1, N'Действия над объектами учета'),
       (2, N'Действия над этикеткой'),
       (3, N'Цвет продукции'),
       (4, N'Принтеры'),
       (5, N'Действия для заявок'),
       (6, N'Приоритеты'),
       (7, N'Статусы заявок'),
       (8, N'Компьютеры'),
       (9, N'Участки упаковки'),
       (10, N'Участки хранения'),
       (11, N'Объекты учёта'),
       (12, N'Назначение при списании'),
       (13, N'Комментарии к пакет-поддонам'),
       (14, N'Размеры этикеток'),
       (15, N'Типы документов');
SET IDENTITY_INSERT handbook OFF;
END
GO

exec dbo.fgw_insert_handbook;

-- Таблица сотрудники.
CREATE TABLE dbo.employee
(
    id_employee    UNIQUEIDENTIFIER PRIMARY KEY,     -- автоматическая генерация uuid
    service_number INT                     NOT NULL, -- табельный номер сотрудника, в дальнейшем логин для входа
    first_name     VARCHAR(50)             NOT NULL, -- имя сотрудника
    last_name      VARCHAR(50)             NOT NULL, -- фамилия сотрудника
    patronymic     VARCHAR(50)             NOT NULL, -- отчество сотрудника
    passwd         VARCHAR(255) DEFAULT '' NOT NULL, -- пароль сотрудника
    role_id        UNIQUEIDENTIFIER        NOT NULL  -- роль сотрудника
);

-- Таблица роли.
CREATE TABLE dbo.role
(
    id_role UNIQUEIDENTIFIER PRIMARY KEY, -- генерирует uuid
    number  INT         NOT NULL,         -- номер роли
    name    VARCHAR(55) NOT NULL          -- название роли
);

-- Таблица вариантов упаковки.
CREATE TABLE dbo.packVariant
(
    id_pack_variant     INT IDENTITY PRIMARY KEY,           -- генерирует и итерирует serial
    prod_id             INT          DEFAULT 0    NOT NULL, -- dbo.catalog.id == prod_id конструкторское наименование продукции (dbo.catalog.handbook_id = 0)
    article             VARCHAR(5)   DEFAULT ''   NOT NULL, -- артикул продукции
    pack_name           VARCHAR(255) DEFAULT ''   NOT NULL, -- наименование продукции на этикетке
    color               INT          DEFAULT 0    NOT NULL, -- цвет продукции (dbo.catalog.handbook_id = 3)
    gl                  INT          DEFAULT 70   NOT NULL, -- gl - цифры, петля Мёбиуса (значит продукцию можно перерабатывать)
    quantity_rows       INT          DEFAULT 0    NOT NULL, -- количество рядов в паллет-поддоне
    quantity_per_rows   INT          DEFAULT 0    NOT NULL, -- количество в ряду в паллет-поддоне
    weight              INT          DEFAULT 0    NOT NULL, -- вес паллет-поддона
    depth               INT          DEFAULT 1000 NOT NULL, -- глубина в мм (стандартно 1000 или 800)
    width               INT          DEFAULT 1200 NOT NULL, -- ширина (стандартно 1200)
    height              INT          DEFAULT 0    NOT NULL, -- высота в мм
    is_food             BIT          DEFAULT 0    NOT NULL, -- пищевая продукция 0-нет/1-да
    is_afraid_moisture  BIT          DEFAULT 0    NOT NULL, -- боится влаги 0-нет/1-да
    is_afraid_sun       BIT          DEFAULT 0    NOT NULL, -- беречь от солнца 0-нет/1-да
    is_eaс              BIT          DEFAULT 0    NOT NULL, -- знак соответствия EAC (маркируют на каждую единицу продукции) 0-нет/1-да
    is_accounting_batch BIT          DEFAULT 0    NOT NULL, -- учет партии 0-нет/1-да
    method_ship         BIT          DEFAULT 0    NOT NULL, -- способ отгрузки 0-АТ/1-ЖД
    shelf_life_months   INT          DEFAULT 0    NOT NULL, -- срок годности в месяцах
    bath_furnace        INT          DEFAULT 0    NOT NULL, -- норме ванной печи
    machine_line        INT          DEFAULT 0    NOT NULL, -- номер машинной линии
    is_manufactured     BIT          DEFAULT 0    NOT NULL, -- изготавливается (производится)
    current_date_batch  DATETIME                  NOT NULL, -- текущая дата партии
    numbering_batch     INT          DEFAULT 1    NOT NULL, -- нумерация партии 0 - автоматическая, 1 - ручная, 2 - с указанной даты
    is_archive          BIT          DEFAULT 0    NOT NULL, -- в архиве
    owner_user          UNIQUEIDENTIFIER          NOT NULL, -- uuid владельца записи
    owner_user_datetime DATETIME                  NOT NULL, -- дата и время записи владельца
    last_user           UNIQUEIDENTIFIER          NOT NULL, -- uuid последнего
    last_user_datetime  DATETIME                  NOT NULL  -- дата и время последней модификации
);

-- Таблица операции над ГП.
CREATE TABLE dbo.operation
(
    id_operation        INT IDENTITY PRIMARY KEY,  -- генерирует и итерирует serial
    type_operation      INT              NOT NULL, -- тип операции (0 - приход, 1 - перемещение, 2 - списание, 3 - продажа, 4 - инвентаризация)
    date_operation      DATETIME         NOT NULL, -- дата создания операции
    created_by_employee INT              NOT NULL, -- табельный номер сотрудника, создавшего операцию
    date_order          DATETIME         NOT NULL, -- дата ордера
    ordered_by_employee INT              NOT NULL, -- табельный номер сотрудника, сформировавшего ордера
    code_accounting_obj INT DEFAULT 0    NOT NULL, -- код объекта учета (0 - паллет-поддон, 1 - форма-комплект)
    appoint             INT DEFAULT 0    NOT NULL, -- назначение при списании (0 - в бой, 1 - на переупаковку)
    owner_user          UNIQUEIDENTIFIER NOT NULL, -- uuid владельца записи
    owner_user_datetime DATETIME         NOT NULL, -- дата и время записи владельца
    last_user           UNIQUEIDENTIFIER NOT NULL, -- uuid последнего
    last_user_datetime  DATETIME         NOT NULL  -- дата и время последней модификации
);

-- Таблица спецификация операций.
CREATE TABLE dbo.operationSpecification
(
    id_operation_specification INT IDENTITY PRIMARY KEY,  -- генерирует и итерирует serial
    operation_id               INT DEFAULT 0    NOT NULL, -- ид операции
    code_accounting_obj        INT DEFAULT 0    NOT NULL, -- код объекта учета (0 - паллет-поддон, 1 - форма-комплект)
    obj_id                     INT DEFAULT 0    NOT NULL, -- ид объект учета obj_id == id_pack_variant
    production_date            DATETIME         NOT NULL, -- дата производства
    storage_from               INT DEFAULT 0    NOT NULL, -- ид участка ОТКУДА
    storage_to                 INT DEFAULT 0    NOT NULL, -- ид участка КУДА
    quantity                   INT DEFAULT 0    NOT NULL, -- количество продукции
    owner_user                 UNIQUEIDENTIFIER NOT NULL, -- uuid владельца записи
    owner_user_datetime        DATETIME         NOT NULL, -- дата и время записи владельца
    last_user                  UNIQUEIDENTIFIER NOT NULL, -- uuid последнего
    last_user_datetime         DATETIME         NOT NULL  -- дата и время последней модификации
);

-- Таблица ордеров.
CREATE TABLE dbo.orderOperation
(
    id_order_operation         INT IDENTITY PRIMARY KEY,
    operation_id               INT DEFAULT 0    NOT NULL, -- ид операции
    operation_specification_id INT DEFAULT 0    NOT NUll, -- ид спецификации операции
    storage_id                 INT DEFAULT 0    NOT NULL, -- catalog.handbook_id == 10 ид склада
    code_accounting_obj        INT DEFAULT 0    NOT NULL, -- код объекта учета (0 - паллет-поддон, 1 - форма-комплект)
    obj_id                     INT DEFAULT 0    NOT NULL, -- ид объект учета obj_id == id_pack_variant
    production_date            DATETIME         NOT NULL, -- дата производства
    order_date                 DATETIME         NOT NULL, -- дата ордера
    quantity                   INT DEFAULT 0    NOT NULL, -- количество
    balance                    INT DEFAULT 0    NOT NULL, -- остаток
    owner_user                 UNIQUEIDENTIFIER NOT NULL, -- uuid владельца записи
    owner_user_datetime        DATETIME         NOT NULL, -- дата и время записи владельца
    last_user                  UNIQUEIDENTIFIER NOT NULL, -- uuid последнего
    last_user_datetime         DATETIME         NOT NULL  -- дата и время последней модификации
);

-- Таблица этикеток
CREATE TABLE dbo.ticket
(
    id_ticket        INT IDENTITY PRIMARY KEY,
    pack_variant_id  INT         DEFAULT 0  NOT NULL, -- ид варианта упаковки
    barcode          VARCHAR(13) DEFAULT '' NOT NULL, -- бар-код
    action_last      INT         DEFAULT 0  NOT NULL, -- последнее действие dbo.actionTicket.action = last_action
    action_last_date DATETIME               NOT NULL, --дата последнего действия dbo.actionTicket.action_date = last_date_action
    production_date  DATETIME               NOT NULL, -- дата производства
    is_repack        BIT         DEFAULT 0  NOT NULL  -- переупакованная (0 - нет, 1 - да)
);

-- Таблица действие над этикетками.
CREATE TABLE dbo.actionTicket
(
    id_action_ticker    INT IDENTITY PRIMARY KEY,
    ticket_id           INT DEFAULT 0 NOT NULL, -- ид этикетки
    action_type         INT DEFAULT 0 NOT NULL, -- catalog.handbook_id == 2 тип действия (1 - печать, 2 - упаковка, 3 - разупаковка, 4 - отгрузка, 5 - оприходование)
    created_by_employee INT           NOT NULL, -- dbo.operation.created_by_employee == operation.created_by_employee табельный номер сотрудника, создавшего операцию
    order_operation_id  INT DEFAULT 0 NOT NULL, -- ид ордера
    last_pack_id        INT DEFAULT 0 NOT NUll  -- на каком участке был упакован
);

CREATE PROCEDURE dbo.fgw_role_all -- ХП возвращает список ролей
    AS
BEGIN
    SET NOCOUNT ON;
SELECT id_role, number, name FROM role;
END
GO

CREATE PROCEDURE dbo.fgw_role_find_by_id -- ХП ищет роль по ИД
    @idRole UNIQUEIDENTIFIER -- ид роль
AS
BEGIN
    SET NOCOUNT ON;
SELECT id_role, number, name FROM role WHERE id_role = @idRole;
END
GO

CREATE PROCEDURE dbo.fgw_role_add -- ХП добавляет роль
    @idRole UNIQUEIDENTIFIER, -- ид роль
    @number INT, -- номер роли
    @name VARCHAR(55) -- наименование роли
AS
BEGIN
    SET NOCOUNT ON;
INSERT INTO role(id_role, number, name) VALUES (@idRole, @number, @name);
END
GO

CREATE PROCEDURE dbo.fgw_role_update -- ХП обновляет роль
    @idRole UNIQUEIDENTIFIER, -- ид роль
    @number INT, -- номер роли
    @name VARCHAR(55) -- наименование роли
AS
BEGIN
    SET NOCOUNT ON;
UPDATE role
SET number = @number,
    name   = @name
WHERE id_role = @idRole;
END
GO

CREATE PROCEDURE dbo.fgw_role_delete_by_id -- ХП удаляет роль по ИД
    @idRole UNIQUEIDENTIFIER -- ид роль
AS
BEGIN
    SET NOCOUNT ON;
    DELETE role WHERE id_role = @idRole;
END
GO

CREATE PROCEDURE dbo.fgw_role_exist -- ХП проверяет на существование роли
    @idRole UNIQUEIDENTIFIER
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (SELECT 1 FROM role WHERE id_role = @idRole)
SELECT CAST(1 AS bit) AS ExistsFlag;
ELSE
SELECT CAST(0 AS bit) AS ExistsFlag;
END
GO

CREATE PROCEDURE dbo.fgw_employee_all -- ХП возвращает список сотрудников
    AS
BEGIN
    SET NOCOUNT ON;
SELECT id_employee, service_number, first_name, last_name, patronymic, passwd, role_id FROM employee;
END
GO

CREATE PROCEDURE dbo.fgw_employee_find_by_id -- ХП ищет сотрудника по ИД
    @idEmployee UNIQUEIDENTIFIER -- ид сотрудника
AS
BEGIN
    SET NOCOUNT ON;
SELECT id_employee, service_number, first_name, last_name, patronymic, passwd, role_id
FROM employee
WHERE id_employee = @idEmployee;
END
GO

CREATE PROCEDURE dbo.fgw_employee_add -- ХП добавляет сотрудника
    @idEmployee UNIQUEIDENTIFIER, -- ид сотрудника
    @serviceNumber INT, -- табельный номер сотрудника, в дальнейшем логин для входа
    @firstName VARCHAR(50), -- имя сотрудника
    @lastName VARCHAR(50), -- фамилия сотрудника
    @patronymic VARCHAR(50),-- отчество сотрудника
    @passwd VARCHAR(255), -- пароль сотрудника
    @roleId UNIQUEIDENTIFIER -- роль сотрудника
AS
BEGIN
    SET NOCOUNT ON;
INSERT INTO employee(id_employee, service_number, first_name, last_name, patronymic, passwd, role_id)
VALUES (@idEmployee,
        @serviceNumber,
        @firstName,
        @lastName,
        @patronymic,
        @passwd,
        @roleId);
END
GO

CREATE PROCEDURE dbo.fgw_employee_update -- ХП обновляет сотрудника
    @idEmployee UNIQUEIDENTIFIER, -- ид сотрудника
    @serviceNumber INT, -- табельный номер сотрудника, в дальнейшем логин для входа
    @firstName VARCHAR(50), -- имя сотрудника
    @lastName VARCHAR(50), -- фамилия сотрудника
    @patronymic VARCHAR(50),-- отчество сотрудника
    @passwd VARCHAR(255), -- пароль сотрудника
    @roleId UNIQUEIDENTIFIER -- роль сотрудника
AS
BEGIN
    SET NOCOUNT ON;
UPDATE employee
SET service_number = @serviceNumber,
    first_name     = @firstName,
    last_name      = @lastName,
    patronymic     = @patronymic,
    passwd         = @passwd,
    role_id        = @roleId
WHERE id_employee = @idEmployee;
END
GO

CREATE PROCEDURE dbo.fgw_employee_delete_by_id -- ХП удаляет сотрудника по ИД
    @idEmployee UNIQUEIDENTIFIER -- ид сотрудника
AS
BEGIN
    SET NOCOUNT ON;
    DELETE employee WHERE id_employee = @idEmployee;
END
GO

CREATE PROCEDURE dbo.fgw_employee_exist -- ХП проверяет на существование сотрудника
    @idEmployee UNIQUEIDENTIFIER
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (SELECT 1 FROM employee WHERE id_employee = @idEmployee)
SELECT CAST(1 AS bit) AS ExistsFlag;
ELSE
SELECT CAST(0 AS bit) AS ExistsFlag;
END
GO
--
CREATE PROCEDURE dbo.fgw_handbook_all -- ХП возвращает список справочников
    AS
BEGIN
    SET NOCOUNT ON;
SELECT id_handbook, name FROM handbook;
END
GO

CREATE PROCEDURE dbo.fgw_handbook_find_by_id -- ХП ищет справочник по ИД
    @idHandbook INT -- ид справочника
AS
BEGIN
    SET NOCOUNT ON;
SELECT id_handbook, name FROM handbook WHERE id_handbook = @idHandbook;
END
GO

CREATE PROCEDURE dbo.fgw_handbook_add -- ХП добавляет справочник
    @name VARCHAR(150) -- наименование справочника
AS
BEGIN
    SET NOCOUNT ON;
INSERT INTO handbook(name) VALUES (@name);
END
GO

CREATE PROCEDURE dbo.fgw_handbook_update -- ХП обновляет справочник
    @idHandbook INT, -- ид справочника
    @name VARCHAR(55) -- наименование справочника
AS
BEGIN
    SET NOCOUNT ON;
UPDATE handbook
SET name = @name
WHERE id_handbook = @idHandbook;
END
GO

CREATE PROCEDURE dbo.fgw_handbook_delete_by_id -- ХП удаляет справочник по ИД
    @idHandbook INT -- ид справочника
AS
BEGIN
    SET NOCOUNT ON;
    DELETE handbook WHERE id_handbook = @idHandbook;
END
GO

CREATE PROCEDURE dbo.fgw_handbook_exist -- ХП проверяет на существование справочника
    @idHandbook INT
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (SELECT 1 FROM handbook WHERE id_handbook = @idHandbook)
SELECT CAST(1 AS bit) AS ExistsFlag;
ELSE
SELECT CAST(0 AS bit) AS ExistsFlag;
END
GO
--


CREATE PROCEDURE dbo.fgw_pack_variant_all -- ХП возвращает список вариантов упаковки
    AS
BEGIN
SELECT id_pack_variant,
       prod_id,
       article,
       pack_name,
       color,
       gl,
       quantity_rows,
       quantity_per_rows,
       weight,
       depth,
       width,
       height,
       is_food,
       is_afraid_moisture,
       is_afraid_sun,
       is_eaс,
       is_accounting_batch,
       method_ship,
       shelf_life_months,
       bath_furnace,
       machine_line,
       is_manufactured,
       current_date_batch,
       numbering_batch,
       is_archive,
       owner_user,
       owner_user_datetime,
       last_user,
       last_user_datetime
FROM packVariant;
END
GO

CREATE PROCEDURE dbo.fgw_catalog_all -- ХП возвращает список каталогов
    AS
BEGIN
SELECT id_catalog,
       parent_id,
       handbook_id,
       record_index,
       name,
       comment,
       handbook_value_int_1,
       handbook_value_int_2,
       handbook_value_decimal_1,
       handbook_value_decimal_2,
       handbook_value_bool_1,
       handbook_value_bool_2,
       is_archive,
       owner_user,
       owner_user_datetime,
       last_user,
       last_user_datetime
FROM catalog;
END
GO

CREATE PROCEDURE dbo.fgw_catalog_add -- ХП добавляет объекты в справочник
    @parentId INT, -- для [id_catalog] родительской записи
    @handbookId INT, -- номер справочника
    @recordIndex INT, -- индекс записи (может повторяться)
    @name VARCHAR(255), -- название
    @comment VARCHAR(8000), -- комментарий
    @handbookValueInt1 INT, -- для [handbook_id] (0 - срок годности в месяцах, 1 - объект учета, 3, 4, 9, 10) дополнительное поле 1 справочника (int)
    @handbookValueInt2 INT, -- дополнительное поле 2 справочника (int)
    @handbookValueDecimal1 DECIMAL(10, 2), -- для [handbook_id] (10 - возможный процент использования) дополнительное поле 1 справочника (decimal)
    @handbookValueDecimal2 DECIMAL(10, 2), -- для [handbook_id] (10 - вместимость, S, V) дополнительное поле 2 справочника (decimal)
    @handbookValueBool1 BIT, -- для [handbook_id] (9 - переупаковка да/нет, 10 - наличие ЖД путей да/нет9 - переупаковка да/нет, 10 - наличие ЖД путей да/нет) дополнительное поле 1 справочника (boolean)
    @handbookValueBool2 BIT, -- дополнительное поле 2 справочника (boolean)
    @isArchive BIT, -- архивная запись да/нет
    @ownerUser UNIQUEIDENTIFIER, -- uuid владельца записи
    @ownerUserDateTime DATETIME, -- дата и время записи владельца
    @lastUser UNIQUEIDENTIFIER, -- uuid последнего
    @lastUserDateTime DATETIME -- дата и время последней модификации
AS
BEGIN
    SET NOCOUNT ON;
INSERT INTO catalog(parent_id, handbook_id, record_index, name, comment, handbook_value_int_1, handbook_value_int_2,
                    handbook_value_decimal_1, handbook_value_decimal_2, handbook_value_bool_1,
                    handbook_value_bool_2, is_archive, owner_user, owner_user_datetime, last_user,
                    last_user_datetime)
VALUES (@parentId,
        @handbookId,
        @recordIndex,
        @name,
        @comment,
        @handbookValueInt1,
        @handbookValueInt2,
        @handbookValueDecimal1,
        @handbookValueDecimal2,
        @handbookValueBool1,
        @handbookValueBool2,
        @isArchive,
        @ownerUser,
        @ownerUserDateTime,
        @lastUser,
        @lastUserDateTime);
END
GO
