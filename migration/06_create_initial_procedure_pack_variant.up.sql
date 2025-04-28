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

CREATE PROCEDURE dbo.fgw_pack_variant_add -- ХП добавляет вариант упаковки продукции
    @ProdId INT = 0, -- dbo.catalog.id == prod_id конструкторское наименование продукции (dbo.catalog.handbook_id = 0)
    @Article VARCHAR(5) = '', -- артикул продукции
    @PackName VARCHAR(255) = '', -- наименование продукции на этикетке
    @Color INT = 0, -- цвет продукции (dbo.catalog.handbook_id = 3)
    @GL INT = 0, -- gl - цифры, петля Мёбиуса (значит продукцию можно перерабатывать 70-79)
    @QuantityRows INT = 0, -- количество рядов в паллет-поддоне
    @QuantityPerRows INT = 0, -- количество в ряду в паллет-поддоне
    @Weight INT = 0, -- вес паллет-поддона
    @Depth INT = 0, -- глубина в мм (стандартно 1000 или 800)
    @Width INT = 0, -- ширина (стандартно 1200)
    @Height INT = 0, -- высота в мм
    @IsFood BIT = 0, -- пищевая продукция 0-нет/1-да
    @IsAfraidMoisture BIT = 0, -- боится влаги 0-нет/1-да
    @IsAfraidSun BIT = 0, -- беречь от солнца 0-нет/1-да
    @IsEAC BIT = 0, -- знак соответствия EAC (маркируют на каждую единицу продукции) 0-нет/1-да
    @IsAccountingBatch BIT = 0, -- учет партии 0-нет/1-да
    @MethodShip BIT = 0, -- способ отгрузки 0-АТ/1-ЖД
    @ShelfLifeMonths INT = 0, -- срок годности в месяцах
    @BathFurnace INT = 0, -- норме ванной печи
    @MachineLine INT = 0, -- номер машинной линии
    @IsManufactured BIT = 0, -- изготавливается (производится)
    @CurrentDateBatch DATETIME = 0, -- текущая дата партии
    @NumberingBatch INT = 0, -- нумерация партии 0 - автоматическая, 1 - ручная, 2 - с указанной даты
    @IsArchive BIT = 0, -- в архиве
    @OwnerUserId UNIQUEIDENTIFIER, -- uuid владельца записи
    @OwnerUserDataTime DATETIME, -- дата и время записи владельца
    @LastUser UNIQUEIDENTIFIER, -- uuid последнего
    @LastUserDateTime DATETIME -- дата и время последней модификации
AS
BEGIN
    SET NOCOUNT ON;
    INSERT INTO packVariant(prod_id, article, pack_name, color, gl, quantity_rows, quantity_per_rows, weight, depth,
                            width, height, is_food, is_afraid_moisture, is_afraid_sun, is_eaс, is_accounting_batch,
                            method_ship, shelf_life_months, bath_furnace, machine_line, is_manufactured,
                            current_date_batch, numbering_batch, is_archive, owner_user, owner_user_datetime, last_user,
                            last_user_datetime)
    VALUES (@ProdId,
            @Article,
            @PackName,
            @Color,
            @GL,
            @QuantityRows,
            @QuantityPerRows,
            @Weight,
            @Depth,
            @Width,
            @Height,
            @IsFood,
            @IsAfraidMoisture,
            @IsAfraidSun,
            @IsEAC,
            @IsAccountingBatch,
            @MethodShip,
            @ShelfLifeMonths,
            @BathFurnace,
            @MachineLine,
            @IsManufactured,
            @CurrentDateBatch,
            @NumberingBatch,
            @IsArchive,
            @OwnerUserId,
            @OwnerUserDataTime,
            @LastUser,
            @LastUserDateTime);
END
GO

CREATE PROCEDURE dbo.fgw_pack_variant_find_by_id -- ХП ищет вариант упаковки по ИД
    @idPackVariant int -- ИД варианта упаковки
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
    FROM packVariant
    WHERE id_pack_variant = @idPackVariant
END
GO

CREATE PROCEDURE dbo.fgw_pack_variant_update -- ХП обновляет объект в варианта упаковки
    @idPackVariant INT,
    @ProdId INT = 0, -- dbo.catalog.id == prod_id конструкторское наименование продукции (dbo.catalog.handbook_id = 0)
    @Article VARCHAR(5) = '', -- артикул продукции
    @PackName VARCHAR(255) = '', -- наименование продукции на этикетке
    @Color INT = 0, -- цвет продукции (dbo.catalog.handbook_id = 3)
    @GL INT = 0, -- gl - цифры, петля Мёбиуса (значит продукцию можно перерабатывать 70-79)
    @QuantityRows INT = 0, -- количество рядов в паллет-поддоне
    @QuantityPerRows INT = 0, -- количество в ряду в паллет-поддоне
    @Weight INT = 0, -- вес паллет-поддона
    @Depth INT = 0, -- глубина в мм (стандартно 1000 или 800)
    @Width INT = 0, -- ширина (стандартно 1200)
    @Height INT = 0, -- высота в мм
    @IsFood BIT = 0, -- пищевая продукция 0-нет/1-да
    @IsAfraidMoisture BIT = 0, -- боится влаги 0-нет/1-да
    @IsAfraidSun BIT = 0, -- беречь от солнца 0-нет/1-да
    @IsEAC BIT = 0, -- знак соответствия EAC (маркируют на каждую единицу продукции) 0-нет/1-да
    @IsAccountingBatch BIT = 0, -- учет партии 0-нет/1-да
    @MethodShip BIT = 0, -- способ отгрузки 0-АТ/1-ЖД
    @ShelfLifeMonths INT = 0, -- срок годности в месяцах
    @BathFurnace INT = 0, -- норме ванной печи
    @MachineLine INT = 0, -- номер машинной линии
    @IsManufactured BIT = 0, -- изготавливается (производится)
    @CurrentDateBatch DATETIME = 0, -- текущая дата партии
    @NumberingBatch INT = 0, -- нумерация партии 0 - автоматическая, 1 - ручная, 2 - с указанной даты
    @IsArchive BIT = 0, -- в архиве
    @OwnerUserId UNIQUEIDENTIFIER, -- uuid владельца записи
    @OwnerUserDataTime DATETIME, -- дата и время записи владельца
    @LastUser UNIQUEIDENTIFIER, -- uuid последнего
    @LastUserDateTime DATETIME -- дата и время последней модификации
AS
BEGIN
    UPDATE packVariant
    SET prod_id             = @ProdId,
        article             = @Article,
        pack_name           = @PackName,
        color               = @Color,
        gl                  = @GL,
        quantity_rows       = @QuantityRows,
        quantity_per_rows   = @QuantityPerRows,
        weight              = @Weight,
        depth               = @Depth,
        width               = @Width,
        height              = @Height,
        is_food             = @IsFood,
        is_afraid_moisture  = @IsAfraidMoisture,
        is_afraid_sun       = @IsAfraidSun,
        is_eaс              = @IsEAC,
        is_accounting_batch = @IsAccountingBatch,
        method_ship         = @MethodShip,
        shelf_life_months   = @ShelfLifeMonths,
        bath_furnace        = @BathFurnace,
        machine_line        = @MachineLine,
        is_manufactured     = @IsManufactured,
        current_date_batch  = @CurrentDateBatch,
        numbering_batch     = @NumberingBatch,
        is_archive          = @IsArchive,
        owner_user          = @OwnerUserId,
        owner_user_datetime = @OwnerUserDataTime,
        last_user           = @LastUser,
        last_user_datetime  = @LastUserDateTime
    WHERE id_pack_variant = @idPackVariant
END
GO

CREATE PROCEDURE dbo.fgw_pack_variant_delete_by_id -- ХП удаляет объект по ИД
    @idPackVariant INT
AS
BEGIN
    DELETE FROM packVariant WHERE id_pack_variant = @idPackVariant;
END
GO