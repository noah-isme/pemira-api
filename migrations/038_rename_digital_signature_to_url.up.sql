-- Rename column to indicate it stores URL, not blob
-- For local DB which might already have digital_signature column
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_schema = 'public' 
               AND table_name = 'voter_status' 
               AND column_name = 'digital_signature') THEN
        ALTER TABLE public.voter_status 
            RENAME COLUMN digital_signature TO digital_signature_url;
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                      WHERE table_schema = 'public' 
                      AND table_name = 'voter_status' 
                      AND column_name = 'digital_signature_url') THEN
        -- Column doesn't exist at all, add it
        ALTER TABLE public.voter_status 
            ADD COLUMN digital_signature_url TEXT;
    END IF;
END $$;

COMMENT ON COLUMN public.voter_status.digital_signature_url IS 'URL to signature image in Supabase Storage';
