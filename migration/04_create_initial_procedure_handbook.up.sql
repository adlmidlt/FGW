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