DELIMITER $$
CREATE PROCEDURE InputRandomRows(IN NumRows BIGINT, IN BatchSize INT)
BEGIN
    DECLARE i INT DEFAULT 0;

PREPARE stmt
    FROM 'INSERT INTO `user` (username, nickname, password, profile_picture)
          VALUES(?, ?, ?, ?, ?)';
SET @v1 = 'user', @v2 = 'nicknameTest', @v = '4p455w0rd', @v4 = 'Pict_user';

START TRANSACTION;
WHILE i < NumRows DO
        SET i = i + 1;
        EXECUTE stmt USING CONCAT(@v1+i), @v2, @v3, CONCAT(@v4+i);
        IF i % BatchSize = 0 THEN
            COMMIT;
START TRANSACTION;
END IF;
END WHILE;
COMMIT;
DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;