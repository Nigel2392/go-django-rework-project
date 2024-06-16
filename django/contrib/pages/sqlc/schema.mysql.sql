CREATE TABLE IF NOT EXISTS PageNode (
    id              BIGINT  AUTO_INCREMENT PRIMARY KEY,
    title           TEXT    NOT NULL,
    path            TEXT    UNIQUE         NOT NULL,
    depth           BIGINT  NOT NULL,
    numchild        BIGINT  NOT NULL,
    status_flags    BIGINT  NOT NULL,
    page_id         BIGINT  NOT NULL,
    typeHash        TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS PageNode_path ON PageNode(path);
CREATE INDEX IF NOT EXISTS PageNode_page_id ON PageNode(page_id);
CREATE INDEX IF NOT EXISTS PageNode_typeHash ON PageNode(typeHash);