CREATE PROCEDURE dbo.fgw_operation_all -- ХП возвращает список операций
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id_operation,
           type_operation,
           create_date,
           created_by_employee,
           date_order,
           closed_by_employee,
           code_accounting_obj,
           appoint,
           owner_user,
           owner_user_datetime,
           last_user,
           last_user_datetime
    FROM operation;
END
GO

CREATE PROCEDURE dbo.fgw_operation_add -- ХП добавляет операцию
    @typeOperation INT = 0, -- тип операции (0 - приход, 1 - перемещение, 2 - списание, 3 - продажа, 4 - инвентаризация)
    @createDate DATETIME, -- дата создания операции
    @createdByEmployee INT = 0, -- табельный номер сотрудника, создавшего операцию
    @dateOrder DATETIME, -- дата ордера
    @closedByEmployee INT = 0, -- табельный номер сотрудника, сформировавшего ордера
    @codeAccountingObj INT = 0, -- код объекта учета (0 - паллет-поддон, 1 - форма-комплект)
    @appoint INT = 0, -- назначение при списании (0 - в бой, 1 - на переупаковку)
    @ownerUser UNIQUEIDENTIFIER, -- uuid владельца записи
    @ownerUserDateTime DATETIME, -- дата и время записи владельца
    @lastUser UNIQUEIDENTIFIER, -- uuid последнего
    @lastUserDateTime DATETIME -- дата и время последней модификации

AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO operation(type_operation, create_date, created_by_employee, date_order,
                          closed_by_employee, code_accounting_obj, appoint, owner_user, owner_user_datetime,
                          last_user, last_user_datetime)
    VALUES (@typeOperation, @createDate, @createdByEmployee,
            @dateOrder, @closedByEmployee, @codeAccountingObj,
            @appoint, @ownerUser, @ownerUserDateTime, @lastUser, @lastUserDateTime)
END
GO

CREATE PROCEDURE dbo.fgw_operation_update -- ХП обновляет операцию
    @idOperation INT,
    @typeOperation INT = 0, -- тип операции (0 - приход, 1 - перемещение, 2 - списание, 3 - продажа, 4 - инвентаризация)
    @createDate DATETIME, -- дата создания операции
    @createdByEmployee INT = 0, -- табельный номер сотрудника, создавшего операцию
    @dateOrder DATETIME, -- дата ордера
    @closedByEmployee INT = 0, -- табельный номер сотрудника, сформировавшего ордера
    @codeAccountingObj INT = 0, -- код объекта учета (0 - паллет-поддон, 1 - форма-комплект)
    @appoint INT = 0, -- назначение при списании (0 - в бой, 1 - на переупаковку)
    @lastUser UNIQUEIDENTIFIER, -- uuid последнего
    @lastUserDateTime DATETIME -- дата и время последней модификации

AS
BEGIN
    SET NOCOUNT ON;

    UPDATE operation
    SET type_operation= @typeOperation,
        create_date= @createDate,
        created_by_employee= @createdByEmployee,
        date_order= @dateOrder,
        closed_by_employee= @closedByEmployee,
        code_accounting_obj= @codeAccountingObj,
        appoint= @appoint,
        last_user= @lastUser,
        last_user_datetime = @lastUserDateTime
    WHERE id_operation = @idOperation
END
GO

CREATE PROCEDURE dbo.fgw_operation_find_by_id -- ХП находит операцию по ИД
@idOperation INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT id_operation,
           type_operation,
           create_date,
           created_by_employee,
           date_order,
           closed_by_employee,
           code_accounting_obj,
           appoint,
           owner_user,
           owner_user_datetime,
           last_user,
           last_user_datetime
    FROM operation
    WHERE id_operation = @idOperation
END
GO

CREATE PROCEDURE dbo.fgw_operation_delete_by_id -- ХП удаляет операцию по ИД
@idOperation INT
AS
BEGIN
    SET NOCOUNT ON;

    DELETE
    FROM operation
    WHERE id_operation = @idOperation
END
GO

CREATE PROCEDURE dbo.fgw_operation_exists -- ХП проверяет операцию на существование
@idOperation INT
AS
BEGIN
    SET NOCOUNT ON;

    IF EXISTS (SELECT 1 FROM operation WHERE id_operation = @idOperation)
        SELECT CAST(1 AS bit) AS ExistsFlag;
    ELSE
        SELECT CAST(0 AS bit) AS ExistsFlag;
END
GO

