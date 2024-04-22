alter table urls
    add constraint short_uniq
        unique (short);