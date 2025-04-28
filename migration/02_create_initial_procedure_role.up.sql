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