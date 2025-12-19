-- Rollback: rename back to digital_signature
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_schema = 'public' 
               AND table_name = 'voter_status' 
               AND column_name = 'digital_signature_url') THEN
        ALTER TABLE public.voter_status 
            RENAME COLUMN digital_signature_url TO digital_signature;
    END IF;
END $$;
