CREATE PROCEDURE dbo.fgw_handbook_all -- ХП возвращает список справочников
AS
BEGIN
    SET NOCOUNT ON;
    SELECT id_handbook, name, owner_user, owner_user_datetime, last_user, last_user_datetime FROM handbook;
END
go

CREATE PROCEDURE fgw_handbook_add_zero_obj
AS
BEGIN
    SET IDENTITY_INSERT handbook ON;
    INSERT INTO handbook (id_handbook, name, owner_user, owner_user_datetime, last_user, last_user_datetime)
    VALUES (0, N'Конструкторское наименование', newid(), getdate(), newid(), getdate())
    SET IDENTITY_INSERT handbook OFF;
END
go

CREATE PROCEDURE dbo.fgw_handbook_find_by_id -- ХП ищет справочник по ИД
    @idHandbook INT -- ид справочника
AS
BEGIN
    SET NOCOUNT ON;
    SELECT id_handbook, name, owner_user, owner_user_datetime, last_user, last_user_datetime
    FROM handbook
    WHERE id_handbook = @idHandbook;
END
go

CREATE PROCEDURE dbo.fgw_handbook_add -- ХП добавляет справочник
    @name VARCHAR(150), -- наименование справочника
    @ownerUser UNIQUEIDENTIFIER, -- uuid владельца записи
    @ownerUserDateTime DATETIME, -- дата и время записи владельца
    @lastUser UNIQUEIDENTIFIER, -- uuid последнего
    @lastUserDateTime DATETIME -- дата и время последней модификации
AS
BEGIN
    SET NOCOUNT ON;
    INSERT INTO handbook(name,
                         owner_user,
                         owner_user_datetime,
                         last_user,
                         last_user_datetime)
    VALUES (@name,
            @ownerUser,
            @ownerUserDateTime,
            @lastUser,
            @lastUserDateTime);
END
go

CREATE PROCEDURE dbo.fgw_handbook_update -- ХП обновляет справочник
    @idHandbook INT, -- ид справочника
    @name VARCHAR(55), -- наименование справочника
    @lastUser UNIQUEIDENTIFIER, -- uuid последнего
    @lastUserDateTime DATETIME -- дата и время последней модификации
AS
BEGIN
    SET NOCOUNT ON;
    UPDATE handbook
    SET name               = @name,
        last_user          = @lastUser,
        last_user_datetime = @lastUserDateTime
    WHERE id_handbook = @idHandbook;
END
go



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