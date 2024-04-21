alter table url
    add constraint short_uniq
        unique (short);