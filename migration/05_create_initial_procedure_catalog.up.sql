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
    @parentId INT = 0, -- для [id_catalog] родительской записи
    @handbookId INT = 0, -- номер справочника
    @recordIndex INT = 0, -- индекс записи (может повторяться)
    @name VARCHAR(255) = '', -- название
    @comment VARCHAR(8000) = '', -- комментарий
    @handbookValueInt1 INT = 0, -- для [handbook_id] (0 - срок годности в месяцах, 1 - объект учета, 3, 4, 9, 10) дополнительное поле 1 справочника (int)
    @handbookValueInt2 INT = 0, -- дополнительное поле 2 справочника (int)
    @handbookValueDecimal1 DECIMAL(10, 2) = 0.00, -- для [handbook_id] (10 - возможный процент использования) дополнительное поле 1 справочника (decimal)
    @handbookValueDecimal2 DECIMAL(10, 2) = 0.00, -- для [handbook_id] (10 - вместимость, S, V) дополнительное поле 2 справочника (decimal)
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

CREATE PROCEDURE dbo.fgw_catalog_find_by_id -- ХП ищет каталог по ИД
    @idCatalog INT -- ид справочника
AS
BEGIN
    SET NOCOUNT ON;
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
    FROM catalog
    WHERE id_catalog = @idCatalog;
END
GO

CREATE PROCEDURE dbo.fgw_catalog_all_find_by_number -- ХП ищет каталоги по номеру справочника
    @idHandbook INT -- номер справочника
AS
BEGIN
    SET NOCOUNT ON;
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
    FROM catalog
    WHERE handbook_id = @idHandbook;
END
GO

CREATE PROCEDURE dbo.fgw_catalog_update -- ХП обновляет объект в справочнике
    @idCatalog INT,
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

    UPDATE catalog
    SET parent_id                = @parentId,
        handbook_id              = @handbookId,
        record_index             = @recordIndex,
        name                     = @name,
        comment                  = @comment,
        handbook_value_int_1     = @handbookValueInt1,
        handbook_value_int_2     = @handbookValueInt2,
        handbook_value_decimal_1 = @handbookValueDecimal1,
        handbook_value_decimal_2 = @handbookValueDecimal2,
        handbook_value_bool_1    = @handbookValueBool1,
        handbook_value_bool_2    = @handbookValueBool2,
        is_archive               = @isArchive,
        owner_user               = @ownerUser,
        owner_user_datetime      = @ownerUserDateTime,
        last_user                = @lastUser,
        last_user_datetime       = @lastUserDateTime
    WHERE id_catalog = @idCatalog
END
GO

CREATE PROCEDURE dbo.fgw_catalog_delete_by_id @idCatalog INT
AS
BEGIN
    SET NOCOUNT ON;

    DELETE FROM catalog WHERE id_catalog = @idCatalog;

END
GO

CREATE PROCEDURE dbo.fgw_catalog_exist -- ХП проверяет на существование каталога
    @idCatalog INT
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (SELECT 1 FROM catalog WHERE id_catalog = @idCatalog)
        SELECT CAST(1 AS bit) AS ExistsFlag;
    ELSE
        SELECT CAST(0 AS bit) AS ExistsFlag;
END
GO

CREATE PROCEDURE dbo.fgw_pack_variant_exist -- ХП проверяет на существование варианта упаковки
@idPackVariant INT
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (SELECT 1 FROM packVariant WHERE id_pack_variant = @idPackVariant)
        SELECT CAST(1 AS bit) AS ExistsFlag;
    ELSE
        SELECT CAST(0 AS bit) AS ExistsFlag;
END
GO