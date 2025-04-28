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